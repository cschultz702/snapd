// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package asserts

import (
	"regexp"
)

// Repair holds an repair assertion which allows running repair
// code to fixup broken systems. It can be limited by series and models.
type Repair struct {
	assertionBase

	series        []string
	architectures []string
	models        []string
}

// BrandID returns the brand identifier that signed this assertion.
func (r *Repair) BrandID() string {
	return r.HeaderString("brand-id")
}

// RepairID returns the "id" of the repair. It should be a short string
// that follows a convention like "REPAIR-123". Similar to a CVE there
// should be a public place to look up details about the repair-id
// (e.g. the snapcraft forum).
func (r *Repair) RepairID() string {
	return r.HeaderString("repair-id")
}

// Architectures returns the architectures that this assertions applies to.
func (r *Repair) Architectures() []string {
	return r.architectures
}

// Series returns the series that this assertion is valid for.
func (r *Repair) Series() []string {
	return r.series
}

// Models returns the models that this assertion is valid for.
// It is a list of "brand-id/model-name" strings.
func (r *Repair) Models() []string {
	return r.models
}

// Implement further consistency checks.
func (r *Repair) checkConsistency(db RODatabase, acck *AccountKey) error {
	// Do the cross-checks when this assertion is actually used,
	// i.e. in the future repair code

	return nil
}

// sanity
var _ consistencyChecker = (*Repair)(nil)

// the repair-id can for now be a sequential number starting with 1
var validRepairID = regexp.MustCompile("^[1-9][0-9]*$")

func assembleRepair(assert assertionBase) (Assertion, error) {
	err := checkAuthorityMatchesBrand(&assert)
	if err != nil {
		return nil, err
	}

	if _, err = checkStringMatchesWhat(assert.headers, "repair-id", "header", validRepairID); err != nil {
		return nil, err
	}

	series, err := checkStringList(assert.headers, "series")
	if err != nil {
		return nil, err
	}
	models, err := checkStringList(assert.headers, "models")
	if err != nil {
		return nil, err
	}
	architectures, err := checkStringList(assert.headers, "architectures")
	if err != nil {
		return nil, err
	}

	return &Repair{
		assertionBase: assert,
		series:        series,
		architectures: architectures,
		models:        models,
	}, nil
}
