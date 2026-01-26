int add(int a, int b) { return a + b; }
int divide(int a, int b) { return a / b; }

int main(void) {
    volatile int x = 100;
    return add(x, x) / divide(x, 1);
}
