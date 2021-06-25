# FAQ

This is a Frequently Asked Questions page.  This is also a Not Frequently Asked But We Are Going To Try To Anticipate
Questions That Might Be Asked In The Future So You Hopefully Won't Feel The Need To Ask Them page.  Modifications welcome.

### Why is this application named trainer?

At The Home Depot, we wrote a microservice ecosystem whose job was to interact with ServiceNow.  It is not necessary
nor appropriate to give further details, the only important fact is that ServiceNow can be abbreviated as SNOW.  On a
day when the lead developer of this project was tired and running out of creative ideas, he decided to mix "snow" and
"pokemon".  This lead to "snowkemon".  Nearly all of the microservices in the ecosystem (and there are quite a few) have
been named after Pokemon characters or other things from the Pokemon universe.  An effort was made to make the names of
the microservices have at least a passing resemblance to what they do.  Since a "trainer" is a person who trains up
different pokemon, and this is an integration testing/mocking microservice, it stuck.

Paradoxically, the lead developer is not a pokemon fan and has only played it in passing.  It just seemed a fun thing to
do.

("Pokemon" is a registered trademark of its owner, who is not The Home Depot.)

### Why was trainer created in the first place?

The short short answer is, because we needed it.

The longer answer is:  We needed an integration testing app that could perform certain tasks, and couldn't find an
app that could do what we needed without shelling out a ton of cash (and even then we weren't sure we could find it at any
price.)  We needed a tool that we could use to exercise the microservices we were developing so that we could ensure they
were behaving as required.  Without going into detail, we were trying to accomplish something rather complicated across
several different services, and this tool proved invaluable for ensuring the success of our deployment/cutover.

The deployment/cutover went almost perfectly, btw.  Thanks mostly to trainer, there were no surprises and everything
just worked.

### Why are you inflicting - er I mean giving this to the world?

We like to give back.  It's one of our core values.  Do we need a better reason?

### Why did you make (insert design choice here)

Because it met our needs.  Yes, it can be complicated to configure.  Yes, it is pretty much a very slow, probably
turing complete (I haven't checked this formally but it does have all of the necessary pieces, like conditionals,
variables, etc) state machine.  But at the end of the day, it gave us what we needed and has been of great use.

It may not meet your needs immediately.  We open sourced it so that you can use it, and contribute back, so that it *can* meet
your needs.  But we hope it's a good starting point.  We can think of many ways it could be improved, and we're sure
the broader community can think of things we didn't.  It's a win for everyone.

Making it simpler to use is certainly a worthy goal, and PRs to that regard are gratefully and eagerly accepted.

### Do you like working at The Home Depot?

Why yes, thanks for asking.  Many companies make a show of living their values, but The Home Depot actually does so.
If you're interested in working for us, please feel free to go to https://careers.homedepot.com and select one that is
right for you.  If you are a tech person, our primary technology centers are in Austin, TX or Atlanta, GA.  If you are
not, you're on a github site reading a FAQ, so we're not sure we believe you.  :-)

### 日本語を話しますか

はい。だけど、私たちの日本語はよくないです。英語を話ってください。  :)

