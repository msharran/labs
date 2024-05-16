#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <argp.h>

const char *argp_program_version = "argp-ex2 1.0";
const char *argp_program_bug_address = "<bug-gnu-utils@gnu.org>";

/* Program documentation. */
static char doc[] = "Argp example #2 -- a pretty minimal program using argp";

/* Our argument parser.  The options, parser, and
   args_doc fields are zero because we have neither options or
   arguments; doc and argp_program_bug_address will be
   used in the output for ‘--help’, and the ‘--version’
   option will print out argp_program_version. */
static struct argp argp = { 0, 0, 0, doc };


void help(char *argv[]) {
	fprintf(stderr, "Usage: %s <pattern> <file>\n", argv[0]);
}

int main(int argc, char *argv[]) {
	argp_parse (&argp, argc, argv, 0, 0, 0);
	// if (argc != 3) {
	// 	help(argv);
	// 	return 1;
	// }
	//
	// FILE *fp = fopen(argv[2], "r");
	// if (fp == NULL) {
	// 	perror("fopen");
	// 	return 1;
	// }
	//
	// char *line = NULL;
	// size_t linecap = 0;
	// ssize_t linelen;
	// while ((linelen = getline(&line, &linecap, fp)) > 0) {
	// 	if (strstr(line, argv[1]) != NULL) {
	// 		fwrite(line, linelen, 1, stdout);
	// 	}
	// }

	return 0;
}
