# Gobatch

Go powered engine that offers features from low level opengl abstraction to UI framework. I created this to separate lot of logic from game am working on. 

## Structure

Engine contains lot of packages and most of them are not even that big. Reason for this is simple. I consider lot more convenient to write `key.A` instead of `ggl.KeyA`, some packages like `rgba` (offers lot of colors) are also generated then there are two package for interpolation `lerp` anc `lerpc`, one focuses on floating points other on colors but actual struct names are same. Another example is `particle`, this way you write `particle.System` and not `drw.ParticleSystem`. Whole engine is thus modular. It prefers components over maybe nice slow abstraction and can be more of a backend.

## Dependencies

I have to mention that engine depends two "languages", `goml` and `goss`. Yes they are named after html and css as they resemble them. You don't have to use them, but they make ui lot easier to develop. See this [repo](github.com/jakubDoka/goml) for documentation of the "languages". I also made vscode extensions for syntax highlighting, link to them can be found in readme of mentioned repository. Errors are handles by [sterr](github.com/jakubDoka/sterr), using it directly might be useful when testing things that depend on gobatch.

## Examples

Nothing is better the learn from code so I wrote couple of examples to show off what engine can do for reference. You can find them all [here](https://github.com/jakubDoka/mlok/tree/main/examples).

![raycasting](https://user-images.githubusercontent.com/60517552/111913276-ae94fc80-8a6d-11eb-8ac6-738b8b561e45.png)

![particles](https://user-images.githubusercontent.com/60517552/111914104-fe28f780-8a70-11eb-9955-b4e4a4c29a06.png)

![chat ui](https://user-images.githubusercontent.com/60517552/111915195-52ce7180-8a75-11eb-8ccf-54f52c800427.png)

How to [create a window](https://github.com/jakubDoka/mlok/wiki/First-window) or [draw a sprite](https://github.com/jakubDoka/mlok/wiki/First-sprite) is also documented.

## Documentation

I am going to be absolutely honest, some comments can be outdated. When i was developing ui package (ant its still in progress), i tried multiple different approaches and commented things too early. There is a lot of documentation and i have to clean it up, document new features and so on. Please open an issue if doc is unclear or is missing so i can prioritize things.

## Contribution

When contributing please keep conventions. If you end up naming lot of struct fields with same prefix, extract them into separate struct and embed/add it to the parent (it has no runtime cost and makes code cleaner). Don't be afraid to introduce new package just to make naming nicer (again you notice it by same prefixes on items). Write tests if possible
or add a exhaustive example of feature use. I have to first understand what code does to decide if its reasonable.

## UI bugs

Ui can contain lot of bugs because of how flexible feature it is. It is just hard to test everything. If something behaves inconveniently open the issue and i will 1) fix it or 2) show you a work around if i cannot fix it (that can happen too).
