/*
 *
 * Copyright 2023 puzzlepassstrengthserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	_ "embed"
	"os"
	"strings"

	grpcserver "github.com/dvaumoron/puzzlegrpcserver"
	"github.com/dvaumoron/puzzlepassstrengthserver/passstrengthserver"
	pb "github.com/dvaumoron/puzzlepassstrengthservice"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

//go:embed version.txt
var version string

func main() {
	// should start with this, to benefit from the call to godotenv
	s := grpcserver.Make(passstrengthserver.PassstrengthKey, version)

	defaultPass := os.Getenv("DEFAULT_PASSWORD")
	rules := readRulesConfig(s.Logger)

	pb.RegisterPassstrengthServer(s, passstrengthserver.New(defaultPass, rules))
	s.Start()
}

func readRulesConfig(logger *otelzap.Logger) map[string]string {
	allLang := strings.Split(os.Getenv("AVAILABLE_LOCALES"), ",")
	localizedRules := make(map[string]string, len(allLang))
	for _, lang := range allLang {
		lang = strings.TrimSpace(lang)

		var pathBuilder strings.Builder
		pathBuilder.WriteString("rules/rules_")
		pathBuilder.WriteString(lang)
		pathBuilder.WriteString(".txt")
		content, err := os.ReadFile(pathBuilder.String())
		if err == nil {
			localizedRules[lang] = strings.TrimSpace(string(content))
		} else {
			logger.Error("Failed to load file", zap.Error(err))
		}
	}
	return localizedRules
}
