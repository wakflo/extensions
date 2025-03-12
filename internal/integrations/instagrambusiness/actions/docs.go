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

package actions

import (
	_ "embed"
)

//go:embed post_reel.md
var postReelDocs string

//go:embed detect_courier.md
var detectCourierDocs string

//go:embed get_a_tracking.md
var getATrackingDocs string

//go:embed get_couriers.md
var getCouriersDocs string

//go:embed get_all_trackings.md
var getAllTrackingsDocs string

//go:embed get_user_couriers.md
var getUserCouriersDocs string

//go:embed mark_tracking_as_completed.md
var markTrackingAsCompletedDocs string

//go:embed retrack_a_tracking.md
var retrackATracking string
