charmer
=======

Motivation
----------

Both [Viper](https://github.com/spf13/viper) and [Cobra](https://github.com/spf13/cobra) are great libraries, but when 
it comes to implementations, combining them turns out to be very verbose. This lib provides ability to define bindings
between these two in form of tagged structure. This should cover the most of the simple use cases.

Limitations
-----------

Charmer only supports subset of types Viper and Cobra support.

Usage
-----

See [tests](./charmer_test.go) for example usage.