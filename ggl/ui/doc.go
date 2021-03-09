// Package ui brings robust ui system based on elements. It is deeply inspired by css and html and it uses
// goml for ui building, this allows you to define your game ui almost as easily as making a website. Goml
// and goss (how original names) gives the visuals and you can use Scene.ID() to access elements easily and
// register event listeners on them. there are some predefined elements you can use as Sprite, Text, Button
// and more. But that does not mean yo cantot create your own element with complex behavior. Every Element
// has ths Module, Module is what defines the behavior. You can easily take already defined module, and build
// on top of it. Don't start from scratch when defining one. Use ModuleBase that already implements all methods
// so you can worry only about methods you need. when you are done with your module, you can register it under
// an identifier by makeing a ModuleFactory and adding it to parser. Then you can freely use the module in your
// goml "code" as you would use button. Look into examples/ui see how to setup the ui loop.
package ui
