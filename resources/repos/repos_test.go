package repos

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	repoA = Repo{
		2,
		1,
		0,
		time.Now(),
		"Netflix/abc",
	}
	repoB = Repo{
		0,
		0,
		2,
		time.Now(),
		"Netflix/abc",
	}
	repoC = Repo{
		1,
		2,
		1,
		time.Now(),
		"Netflix/abc",
	}
	repoD = Repo{
		0,
		0,
		1,
		time.Now(),
		"Netflix/zyx",
	}
)

// TestSortByField_Forks is asserting that we sort by forks and name in ascending fashion
func TestSortByField_Forks(t *testing.T) {
	assert := assert.New(t)
	r := Repos{repoA, repoB, repoC, repoD}
	expected := Repos{repoB, repoD, repoC, repoA}
	sortedForks := r.sortByField("forks")
	assert.Equal(expected, sortedForks)
}

// TestSortByField_OpenIssues is asserting that we sort by open_issues and name in ascending fashion
func TestSortByField_OpenIssues(t *testing.T) {
	assert := assert.New(t)
	r := Repos{repoA, repoB, repoC, repoD}
	expected := Repos{repoB, repoD, repoA, repoC}
	sortedForks := r.sortByField("open_issues")
	assert.Equal(expected, sortedForks)
}

// TestSortByField_Stars is asserting that we sort by stars and name in ascending fashion
func TestSortByField_Stars(t *testing.T) {
	assert := assert.New(t)
	r := Repos{repoA, repoB, repoC, repoD}
	expected := Repos{repoA, repoC, repoD, repoB}
	sortedForks := r.sortByField("stars")
	assert.Equal(expected, sortedForks)
}
