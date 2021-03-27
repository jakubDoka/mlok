// Package ui brings robust ui system based on elements. It is deeply inspired by css and html and it uses
// goml for ui building, this allows you to define your game ui almost as easily as making a website. Goml
// and goss (how original names!) gives the visuals and you can use Scene.ID() to access elements easily then
// register event listeners on them. there are some predefined elements you can use as Sprite, Text, Button
// and more. But that does not mean yo cantot create your own element with complex behavior. Every Element
// has ths Module, that is what defines the behavior. You can easily take already defined module, and build
// on top of it. Don't start from scratch when defining one. Use ModuleBase that already implements all methods
// so you can worry only about methods you need. when you are done with your module, you can register it under
// an identifier by makeing a ModuleFactory and adding it to parser. Then you can freely use the module in your
// goml "code" as you would use for example button. Look into examples/ui to see how to setup the ui loop.
//
// Style Docs
//
// When going trough docs you may come across style documentation. Its in following form:
// 	field_name: type //comment
// This notation describes that if field with given type is added to style of element, something will be altered.
// For example:
//	background: rgba // defines the background color of element
// If you provide `background: blue;` in goml, background will be blue. Now lets explane the type.
//
// rgba - color that can be defined as:
//	name 					// comes from rgba package color map
//	float					// white color with altered Alpha channel
//	float float float		// eed green and blue channels (0 - 1)
//	float float float float	// all channels specified
//	hex (ffffff=white)		// hex notation
//
// vec - vector type
//	float		// results into vector with equal dimensions
//	float float	// each dimension can be differrent
//
// aabb - defines collection of four floats coresponding to left bottom right top respectively
//	float float float float	// most verbose declaration
//	float float				// left and right coresponds to first float, top and bottom to second
//	float					// all four floats coresponds to given value
//
// bool - can be 'true' of 'false'
//
//
package ui
