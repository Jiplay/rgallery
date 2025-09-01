package sizes_test

import (
	"html/template"
	"testing"

	"github.com/robbymilo/rgallery/pkg/sizes"
	"github.com/robbymilo/rgallery/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestSrcset(t *testing.T) {
	var c = types.Conf{
		IncludeOriginals: false,
	}
	var hash uint32 = 3884452138
	ans := template.Srcset(`/img/3884452138/200 200w, /img/3884452138/400 400w, /img/3884452138/800 800w, /img/3884452138/1200 1200w, /img/3884452138/1800 1800w, /img/3884452138/2400 2400w, /img/3884452138/3724 3724w`)
	assert.EqualValues(t, ans, sizes.Srcset(hash, 3724, "", c), "they should be equal")
	ans1 := template.Srcset(`/img/3884452138/200 200w, /img/3884452138/400 400w, /img/3884452138/800 800w, /img/3884452138/1200 1200w, /img/3884452138/1800 1800w, /img/3884452138/2400 2400w, /img/3884452138/4000 4000w`)
	assert.EqualValues(t, ans1, sizes.Srcset(hash, 4123, "", c), "they should be equal")
	ans2 := template.Srcset(`/img/3884452138/200 200w, /img/3884452138/300 300w`)
	assert.EqualValues(t, ans2, sizes.Srcset(hash, 300, "", c), "they should be equal")
}
