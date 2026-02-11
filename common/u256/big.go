// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package u256

import (
	"github.com/holiman/uint256"
)

// Common big integers often used
var (
	Num0  = uint256.NewInt(0)
	Num1  = uint256.NewInt(1)
	Num2  = uint256.NewInt(2)
	Num4  = uint256.NewInt(4)
	Num8  = uint256.NewInt(8)
	Num27 = uint256.NewInt(27)
	Num28 = uint256.NewInt(28)
	Num32 = uint256.NewInt(32)
	Num35 = uint256.NewInt(35)
	Num100 = uint256.NewInt(100)

	// Aliases for compatibility with common2/u256
	N0   = Num0
	N1   = Num1
	N2   = Num2
	N4   = Num4
	N8   = Num8
	N27  = Num27
	N28  = Num28
	N32  = Num32
	N35  = Num35
	N100 = Num100
)
