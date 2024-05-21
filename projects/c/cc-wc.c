#include <stdio.h>
#include <string.h>

int get_char_c(FILE *f) {
    int count = 0;
    char c;
    while ((c = fgetc(f)) != EOF) {
        count++;
    }
    return count;
}

int main(int argc, char *argv[]) {
    if (argc < 3) {
        printf("cc-wc: missing args -- got=%d, want=%d\n", argc, 3);
        fprintf(stderr,"usage: wc [-c] [file ...]\n");
        return 1;
    }

    char *flag = argv[1];
    char *file = argv[2];

    FILE *f = fopen(file, "r");
    if (f == NULL) {
        perror("fopen");
        return 1;
    }

    if (strcmp(flag, "-c") == 0) {
        int cc = get_char_c(f);
        printf("      %d %s\n", cc, file);
    }

    int retval = fclose(f);
    if (retval != 0){
        perror("fclose");
        return retval;
    }
    return 0;
}
