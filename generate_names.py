import names

name_set = set()
for _ in range(20000):
    name_set.add(names.get_full_name().split(' ')[0])

names_list = list(name_set)
names_list.sort()

# Save set to file
with open('names.txt', 'w') as f:
    for name in names_list:
        f.write(name + '\n')