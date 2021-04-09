#include <stdio.h>

int main() {
    FILE *fp = NULL;
    fp=fopen("test.text","w+");
    fprintf(fp,"first line\n");
    fprintf(fp,"second line\n");
    fclose(p);
    return 0;
}