/*
 * Copyright (c) 2021 Michael Morris. All Rights Reserved.
 *
 * Licensed under the MIT license (the "License"). You may not use this file except in compliance
 * with the License. A copy of the License is located at
 *
 * https://github.com/mmmorris1975/aws-runas/blob/master/LICENSE
 *
 * or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License
 * for the specific language governing permissions and limitations under the License.
 */

package cli

import (
	"flag"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestUpdateCmd_Action(t *testing.T) {
	ctx := cli.NewContext(App, new(flag.FlagSet), nil)
	if err := updateCmd.Run(ctx); err != nil {
		t.Error(err)
	}
}
