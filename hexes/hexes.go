// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexes

// TribeNet maps are what Redblob Games calls "flat-top" with "offset even-q coordinates."

// That works well for the big map, but to make the math simple for the steps of
// a move, we're going to use "doubleheight coordinates." We'll assume that the
// first hex of each step is (0,0) and go from there.

// When we start placing tiles on the map, we will have to convert them to the
// even-q offset coordinates so that we can label them. We'll also have to do
// the same when we have units teleport across the map.
