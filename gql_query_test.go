package gjxy

import (
	"testing"
)

const q0 = `abced `

const q1 = `{
	HHHhero {
	  name
	}
  }
`
const q2 = `{
	newName: GGGhero {
	  name
	}
  }
`
const q3 = `{
	FFFhuman(id: "1000") {
	  name
	  height
	}
  }
`
const q4 = `{
	empireHero: DDDhero(episode: EMPIRE) {
	  name
	}
	jediHero: EEEhero(episode: JEDI) {
	  name
	}
  }
`
const q5 = `query HeroNameAndFriends {
	CCChero {
	  name
	  friends {
		name
	  }
	}
  }
`
const q6 = `query HeroNameAndFriends($episode: Episode) {
	BBBhero(episode: $episode) {
	  name
	  friends {
		name
	  }
	}
  }
`
const q7 = `query HeroNameAndFriends($episode: Episode) {
	newName: AAAhero(episode: $episode) {
	  name
	  friends {
		name
	  }
	}
  }
`

var Qs = []string{q0, q1, q2, q3, q4, q5, q6, q7}

func TestGet1stObjInQry(t *testing.T) {
	for i := 0; i <= 7; i++ {
		fPln(Get1stObjInQry(Qs[i]))
	}
}
