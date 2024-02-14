package node

import (
	"flag-example/cmd"
	"flag-example/config/params"

	"github.com/urfave/cli/v2"
)

func configureNetwork(cliCtx *cli.Context) {
	// Flag에 BootStrap Node가 등록되어 있으면 해당 정보 등록
	if len(cliCtx.StringSlice(cmd.BootstrapNode.Name)) > 0 {
		c := params.CurieNetworkConfig()
		c.BootstrapNodes = cliCtx.StringSlice(cmd.BootstrapNode.Name)
		params.OverrideCurieNetworkConfig(c)
	}
	// ContractDeploymentBlock이 뭔지 모르겠네,,
}
