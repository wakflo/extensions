// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trackingmore

import (
	"github.com/wakflo/extensions/internal/integrations/trackingmore/actions"
	"github.com/wakflo/extensions/internal/integrations/trackingmore/shared"
	"github.com/wakflo/extensions/internal/integrations/trackingmore/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

type TrackingMore struct{}

func (n *TrackingMore) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.TrackingMoreSharedAuth,
	}
}

func (n *TrackingMore) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewPaymentTrigger(),
	}
}

func (n *TrackingMore) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewTrackAPackageAction(),
		actions.NewCreateTrackingsAction(),
	}
}

func NewTrackingMore() sdk.Integration {
	return &TrackingMore{}
}

var Integration = sdk.Register(NewTrackingMore())
