#include <stdio.h>

main() {
   read_file();
}

int read_file (){
    FILE *fp=NULL;
    char buf[255];

    fp=fopen("test.text","r");
    fscanf(fp,"%s",buf);
    printf("1: %s\n",buf);

}

int write_file() {
     FILE *fp =NULL;
    fp=fopen("test.text","w+");
    fprintf(fp,"first line\n");
    fprintf(fp,"second line\n");
    fclose(fp);
    return 0;
}