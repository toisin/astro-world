**** How does server.server.init() get called? It seems to be called more than once.
Answer:
http://stackoverflow.com/questions/24790175/when-is-the-init-function-in-go-golang-run
(when packages are loaded, all var assignments are called first, then init() are always called, then main() called. There can be zero or more init(). There can only be zero or 1 main())


*** Why can't I return a pointer to an interface
http://openmymind.net/Things-I-Wish-Someone-Had-Told-Me-About-Go/

*** Setting initial state in React
https://facebook.github.io/react/tips/props-in-getInitialState-as-anti-pattern.html