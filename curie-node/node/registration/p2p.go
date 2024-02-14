package registration

import (
	"flag-example/cmd"
	"flag-example/config/params"
	"log"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func P2PPreregistration(cliCtx *cli.Context) (bootstrapNodeAddrs []string, dataDir string, err error) {
	// 부트스트랩 노드 정보 가져오기
	bootnodesTemp := params.CurieNetworkConfig().BootstrapNodes
	bootstrapNodeAddrs = make([]string, 0)
	for _, addr := range bootnodesTemp {
		if len(bootnodesTemp) == 0 {
			continue
		}
		if filepath.Ext(addr) == ".yaml" {
			// flag에 .yaml 파일을 넘겼다면 이것을 읽고 - 이건 생략하기
		} else {
			// flag에 그냥 주소만을 넘겼다면 바로 부트스트랩 노드로 넘기기
			bootstrapNodeAddrs = append(bootstrapNodeAddrs, addr)
		}
	}

	// 기본 디렉토리 지정하기
	dataDir = cliCtx.String(cmd.DataDirFlag.Name)
	if dataDir == "" {
		dataDir = cmd.DefaultDataDir()
		if dataDir == "" {
			log.Fatal(
				"Could not determine your system's HOME path, please specify a --datadir you wish " +
					"to use for your chain data",
			)
		}
	}
	return
}
