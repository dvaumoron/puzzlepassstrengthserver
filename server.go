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
	"log"
	"os"
	"strings"

	grpcserver "github.com/dvaumoron/puzzlegrpcserver"
	"github.com/dvaumoron/puzzlepassstrengthserver/passstrengthserver"
	pb "github.com/dvaumoron/puzzlepassstrengthservice"
)

func main() {
	// should start with this, to benefit from the call to godotenv
	s := grpcserver.New()

	defaultPass := os.Getenv("DEFAULT_PASSWORD")
	rules := readRulesConfig()

	pb.RegisterPassstrengthServer(s, passstrengthserver.New(defaultPass, rules))
	s.Start()
}

func readRulesConfig() map[string]string {
	allLang := strings.Split(os.Getenv("AVAILABLE_LOCALES"), ",")
	localizedRules := make(map[string]string, len(allLang))
	for _, lang := range allLang {
		var pathBuilder strings.Builder
		pathBuilder.WriteString("rules/rules_")
		pathBuilder.WriteString(strings.TrimSpace(lang))
		pathBuilder.WriteString(".txt")
		path := pathBuilder.String()

		content, err := os.ReadFile(path)
		if err == nil {
			localizedRules[lang] = strings.TrimSpace(string(content))
		} else {
			log.Println("Failed to load file :", err)
		}
	}
	return localizedRules
}
