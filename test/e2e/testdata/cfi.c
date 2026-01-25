typedef int (*func_ptr)(int);

int add_one(int x) { return x + 1; }

int main(void) {
    func_ptr fn = add_one;
    return fn(5) > 0 ? 0 : 1;
}
