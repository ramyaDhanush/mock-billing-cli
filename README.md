# mock-billing-cli
A golang application to mock the billing system in super markets

**Features**

1. View all items & items with filter
2. Refill items with admin access
3. Purchase an item
4. Customer Details
5. Admin Details

**Usage**

1. Move to `cmd` directory
2. Run `go mod init` & `go mod tidy` to create `go.mod` & `go.sum` files
3. Run the cmd.exe file with flags  
   i)    view (View items of inventory)  
         -- all, id, category, price-above, price-below  
   ii)   refill (refill stock with admin access)  
   iii)  bill (enter billing process)  
   iv)   customer (view customer list)  
    v)   admin (view admin list)  


