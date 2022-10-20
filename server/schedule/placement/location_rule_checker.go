// Copyright 2021 TiKV Project Authors.
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

package placement

import (
	"strconv"

	"github.com/pingcap/log"
	"github.com/tikv/pd/pkg/codec"
	"go.uber.org/zap"
)

func (m *RuleManager) UpdateLoacationRule(num int64) error {
	m.Lock()
	defer m.Unlock()
	p := m.beginPatch()
	updated := 0
	for key, rule := range p.c.rules {
		if rule.GroupID == "tiflash" && len(rule.LabelConstraints) == 1 && (codec.Key(rule.StartKey).TableID()&0x0100_0000_0000_0000) > 0 {
			newRule := rule.Clone()
			tid := codec.Key(rule.StartKey).TableID()
			bid := (tid & 0x00_FFFF_00_0000_0000) >> 48
			instanceID := bid % num
			newRule.LabelConstraints = append(newRule.LabelConstraints, LabelConstraint{"bucket_id", In, []string{strconv.Itoa(int(instanceID))}})
			p.mut.rules[key] = newRule
			updated++
		}
	}
	if err := m.tryCommitPatch(p); err != nil {
		return err
	}
	log.Info("placement location rules updated", zap.Int64("num", num), zap.Int("updated", updated))
	return nil
}
