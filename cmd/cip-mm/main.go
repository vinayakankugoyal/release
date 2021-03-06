/*
Copyright 2020 The Kubernetes Authors.

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

package main

import (
	"context"
	"flag"

	// TODO: Use k/release/pkg/log instead

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	reg "sigs.k8s.io/k8s-container-image-promoter/pkg/dockerregistry"
)

var cmd = &cobra.Command{
	Short: "cip-mm → Container Image Promoter - Manifest Modificator",
	Long: `cip-mm → Container Image Promoter - Manifest Modificator
	
This tool **m**odifies promoter **m**anifests. For now it dumps some filtered
subset of a staging GCR and merges those contents back into a given promoter
manifest.`,
	Use:           "cip-mm",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(commandLineOpts)
	},
}

type commandLineOptions struct {
	baseDir      string
	stagingRepo  string
	filterImage  string
	filterDigest string
	filterTag    string
}

var commandLineOpts = &commandLineOptions{}

func init() {
	cmd.PersistentFlags().StringVar(
		&commandLineOpts.baseDir,
		"base_dir",
		"",
		"the manifest directory to look at and modify",
	)
	cmd.PersistentFlags().StringVar(
		&commandLineOpts.stagingRepo,
		"staging_repo",
		"",
		"the staging repo which we want to read from",
	)
	cmd.PersistentFlags().StringVar(
		&commandLineOpts.filterImage,
		"filter_image",
		"",
		"filter staging repo by this image name",
	)
	cmd.PersistentFlags().StringVar(
		&commandLineOpts.filterDigest,
		"filter_digest",
		"",
		"filter images by this digest",
	)
	cmd.PersistentFlags().StringVar(
		&commandLineOpts.filterTag,
		"filter_tag",
		"",
		"filter images by this tag",
	)

	// Add the flags from klog
	flagsCompat := &flag.FlagSet{}
	klog.InitFlags(flagsCompat)
	cmd.PersistentFlags().AddGoFlagSet(flagsCompat)
}

func main() {
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}

// nolint[lll]
func run(cmdOpts *commandLineOptions) error {
	opt := reg.GrowManifestOptions{}
	if err := opt.Populate(
		cmdOpts.baseDir, cmdOpts.stagingRepo,
		cmdOpts.filterImage, cmdOpts.filterDigest,
		cmdOpts.filterTag); err != nil {
		return err
	}

	if err := opt.Validate(); err != nil {
		return err
	}

	ctx := context.Background()
	return reg.GrowManifest(ctx, opt)
}
