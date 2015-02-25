package sources

import (
	"github.com/tyba/opensource-search/sources/social/http"

	. "gopkg.in/check.v1"
)

func (s *SourcesSuite) TestGithub_SearchByEmail(c *C) {
	a := NewGithub(http.NewClient(true))
	g, err := a.GetProfileByURL("https://github.com/mcuadros")
	c.Assert(err, IsNil)
	c.Assert(g.Username, Equals, "mcuadros")
	c.Assert(g.FullName, Equals, "Máximo Cuadros")
	c.Assert(g.Location, Equals, "Madrid, Spain")
	c.Assert(g.Email, Equals, "mcuadros@gmail.com")
	c.Assert(g.Description, Equals, "mcuadros has 64 repositories written in PHP, Go, and Shell. Follow their code on GitHub.")
	c.Assert(g.JoinDate.Unix(), Equals, int64(1332676111))
	c.Assert(g.Organizations, HasLen, 4)
	c.Assert(g.Organizations[0], Equals, "/sourcegraph")
	c.Assert(g.Repositories, HasLen, 5)
	c.Assert(g.Repositories[0].Name, Equals, "dockership")
	c.Assert(g.Repositories[0].Url, Equals, "/mcuadros/dockership")
	c.Assert(g.Repositories[0].Owner, Equals, "mcuadros")
	c.Assert(g.Repositories[0].Stars, Not(Equals), 0)
	c.Assert(g.Contributions, HasLen, 5)
	c.Assert(g.Contributions[0].Name, Equals, "mongofill")
	c.Assert(g.Contributions[0].Url, Equals, "/mongofill/mongofill")
	c.Assert(g.Contributions[0].Owner, Equals, "mongofill")
	c.Assert(g.Contributions[0].Stars, Not(Equals), 0)
	c.Assert(g.Followers, Not(Equals), 0)
	c.Assert(g.Starred, Not(Equals), 0)
	c.Assert(g.Following, Not(Equals), 0)
	c.Assert(g.TotalContributions, Not(Equals), 0)

	g, err = a.GetProfileByURL("https://github.com/philips")
	c.Assert(g.Username, Equals, "philips")
	c.Assert(g.WorksFor, Equals, "CoreOS, Inc")
	c.Assert(g.Url, Equals, "https://github.com/philips")
}