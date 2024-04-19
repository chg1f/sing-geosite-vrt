package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"golang.org/x/sync/errgroup"

	"github.com/google/go-github/v45/github"
)

var (
	githubClient *github.Client

	outputDir, _ = filepath.Abs("rule-set")
	generates    []string
)

func init() {
	accessToken, loaded := os.LookupEnv("ACCESS_TOKEN")
	if !loaded {
		githubClient = github.NewClient(nil)
		return
	}
	transport := &github.BasicAuthTransport{
		Username: accessToken,
	}

	githubClient = github.NewClient(transport.Client())
	os.RemoveAll(outputDir)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatal(err)
	}
}

func getLatestRelease(from string) (*github.RepositoryRelease, error) {
	names := strings.SplitN(from, "/", 2)
	latestRelease, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), names[0], names[1])
	if err != nil {
		return nil, err
	}
	return latestRelease, err
}

func fetch(uri *string) ([]byte, error) {
	log.Info("download ", *uri)
	response, err := http.Get(*uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func download(release *github.RepositoryRelease, assetName string) ([]byte, error) {
	asset := common.Find(release.Assets, func(it *github.ReleaseAsset) bool {
		return *it.Name == assetName
	})
	if asset == nil {
		return nil, E.New(assetName+" not found in upstream release ", release.Name)
	}
	data, err := fetch(asset.BrowserDownloadURL)
	if err != nil {
		return nil, err
	}
	checksumAsset := common.Find(release.Assets, func(it *github.ReleaseAsset) bool {
		return *it.Name == assetName+".sha256sum"
	})
	if checksumAsset != nil {
		remoteChecksum, err := fetch(checksumAsset.BrowserDownloadURL)
		if err != nil {
			return nil, err
		}
		checksum := sha256.Sum256(data)
		if hex.EncodeToString(checksum[:]) != string(remoteChecksum[:64]) {
			return nil, E.New("checksum mismatch")
		}
	}
	return data, nil
}

func generateSource(plainRuleSet option.PlainRuleSet, name string) error {
	bs, err := json.MarshalIndent(option.PlainRuleSetCompat{
		Version: 1,
		Options: plainRuleSet,
	}, "", "  ")
	if err != nil {
		return err
	}
	generates = append(generates, name+".json")
	return os.WriteFile(filepath.Join(outputDir, name+".json"), bs, 0o644)
}

func generateBinary(plainRuleSet option.PlainRuleSet, name string) error {
	output, err := os.Create(filepath.Join(outputDir, name+".srs"))
	if err != nil {
		return err
	}
	defer output.Close()
	err = srs.Write(output, plainRuleSet)
	if err != nil {
		return err
	}
	generates = append(generates, name+".srs")
	return nil
}

func setActionOutput(name string, content string) {
	os.Stdout.WriteString("::set-output name=" + name + "::" + content + "\n")
}

func main() {
	var eg errgroup.Group
	eg.Go(func() error {
		sourceRelease, err := getLatestRelease("Loyalsoldier/clash-rules")
		if err != nil {
			return err
		}
		log.Warn("clash-rules from " + *sourceRelease.TagName)
		return generateClashRules(sourceRelease,
			"apple.txt",
			"cncidr.txt",
			"gfw.txt",
			"greatfire.txt",
			"lancidr.txt",
			"proxy.txt",
			"telegramcidr.txt",
			"applications.txt",
			"direct.txt",
			"google.txt",
			"icloud.txt",
			"private.txt",
			"reject.txt",
			"tld-not-cn.txt",
		)
	})
	eg.Go(func() error {
		sourceRelease, err := getLatestRelease("Loyalsoldier/v2ray-rules-dat")
		if err != nil {
			return err
		}
		log.Warn("v2ray-rules-dat from " + *sourceRelease.TagName)
		return generateV2rayRulesDat(sourceRelease,
			"geosite.db",
			"geosite-cn.db",
		)
	})
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	sort.Strings(generates)
	os.WriteFile(filepath.Join(outputDir, ".rule_set.txt"), []byte(strings.Join(generates, "\n")), 0o644)
	setActionOutput("tag", time.Now().Format("20060102150405"))
}
