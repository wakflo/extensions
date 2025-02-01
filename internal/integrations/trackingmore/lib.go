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
	"github.com/wakflo/go-sdk/integration"
)

type TrackingMore struct{}

func (n *TrackingMore) Auth() *integration.Auth {
	return &integration.Auth{
		Required: true,
		Schema:   *shared.TrackingMoreSharedAuth,
	}
}

func (n *TrackingMore) Triggers() []integration.Trigger {
	return []integration.Trigger{
		triggers.NewNewPaymentTrigger(),
	}
}

func (n *TrackingMore) Actions() []integration.Action {
	return []integration.Action{
		actions.NewTrackAPackageAction(),
		actions.NewCreateTrackingsAction(),
	}
}

func NewTrackingMore() integration.Integration {
	return &TrackingMore{}
}

var Integration = integration.Register(NewTrackingMore())
