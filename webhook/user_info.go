/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"context"

	"github.com/knative/pkg/apis"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CreatorAnnotationSuffix is the suffix of the annotation key to describe
	// the user that created the resource.
	CreatorAnnotationSuffix = "/creator"

	// UpdaterAnnotationSuffix is the suffix of the annotation key to describe
	// the user who last modified the resource.
	UpdaterAnnotationSuffix = "/updater"
)

// SetUserInfoAnnotations sets creator and updater annotations on a resource.
func SetUserInfoAnnotations(resource apis.HasSpec, ctx context.Context, groupName string) {
	if ui := apis.GetUserInfo(ctx); ui != nil {
		objectMetaAccessor, ok := resource.(metav1.ObjectMetaAccessor)
		if !ok {
			return
		}

		annotations := objectMetaAccessor.GetObjectMeta().GetAnnotations()
		if annotations == nil {
			annotations = map[string]string{}
			defer objectMetaAccessor.GetObjectMeta().SetAnnotations(annotations)
		}

		if apis.IsInUpdate(ctx) {
			old := apis.GetBaseline(ctx).(apis.HasSpec)
			if equality.Semantic.DeepEqual(old.GetUntypedSpec(), resource.GetUntypedSpec()) {
				return
			}
			annotations[groupName+UpdaterAnnotationSuffix] = ui.Username
		} else {
			annotations[groupName+CreatorAnnotationSuffix] = ui.Username
			annotations[groupName+UpdaterAnnotationSuffix] = ui.Username
		}
	}
}
