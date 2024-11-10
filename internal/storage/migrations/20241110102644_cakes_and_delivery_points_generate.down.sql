DELETE FROM cakes WHERE description IN (
    'Chocolate Cake', 'Vanilla Cake', 'Red Velvet Cake',
    'Lemon Cake', 'Carrot Cake', 'Cheesecake',
    'Black Forest Cake', 'Fruit Cake', 'Pineapple Cake',
    'Coffee Cake', 'Brownie Cake', 'Banana Cake',
    'Strawberry Cake', 'Mango Cake', 'Marble Cake',
    'Tiramisu', 'Red Wine Cake', 'Peanut Butter Cake',
    'Almond Cake', 'Orange Cake', 'Raspberry Cake',
    'Nut Cake', 'Coconut Cake', 'Pistachio Cake',
    'Millefeuille', 'Lavender Cake', 'Matcha Cake',
    'Frozen Yogurt Cake', 'Ginger Cake', 'Zebra Cake',
    'Sponge Cake', 'Cupcake', 'Baklava Cake',
    'Ice Cream Cake'
);

DELETE FROM delivery_points WHERE address IN (
    '123 Main St', '456 Oak St', '789 Pine St',
    '135 Maple Ave', '246 Elm St', '357 Birch Rd',
    '468 Cedar Blvd', '579 Spruce Ct', '680 Walnut St',
    '791 Cherry Ln', '902 Aspen Dr', '1 Oak St',
    '2 Walnut St', '3 Pine St', '4 Maple Ave',
    '5 Elm St', '6 Cedar Blvd', '7 Birch Rd',
    '8 Spruce Ct', '9 Cherry Ln', '10 Aspen Dr',
    '11 Oak St', '12 Walnut St', '13 Pine St',
    '14 Maple Ave'
);
