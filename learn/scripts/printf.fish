# printf uses the format string FORMAT to print the ARGUMENT arguments. This means that it takes format specifiers in the format string and replaces each with an argument.

# The FORMAT argument is re-used as many times as necessary to convert all of the given arguments. So `printf %s\n flounder catfish clownfish shark` will print four lines.
printf %s\n flounder catfish clownfish shark

# Unlike echo, printf does not append a new line unless it is specified as part of the string.

# It doesn’t support any options, so there is no need for a -- separator, which makes it easier to use for arbitrary input than echo. [1]
