# Go Ecommerce Cart

![Structure Of Project](https://i.imgur.com/sToet5f.png)

This project is a Go-based ecommerce cart implementation. It provides a simple and efficient way to manage a shopping cart for an ecommerce website.

## Features

- Add products to the cart
- Remove products from the cart
- Update the quantity of products in the cart
- Calculate the total price of the cart
- Apply discounts and promotions
- Store cart data persistently

## Installation

1. Clone the repository:

    ```shell
    git clone https://github.com/your-username/go-ecommerce-cart.git
    ```

2. Navigate to the project directory:

    ```shell
    cd go-ecommerce-cart
    ```

3. Install the dependencies:

    ```shell
    go mod download
    ```

4. Build the project:

    ```shell
    go build
    ```

## Usage

1. Import the package in your Go code:

    ```go
    import "github.com/your-username/go-ecommerce-cart"
    ```

2. Create a new cart instance:

    ```go
    cart := cart.NewCart()
    ```

3. Add products to the cart:

    ```go
    cart.AddProduct(product)
    ```

4. Remove products from the cart:

    ```go
    cart.RemoveProduct(productID)
    ```

5. Update the quantity of products in the cart:

    ```go
    cart.UpdateQuantity(productID, quantity)
    ```

6. Calculate the total price of the cart:

    ```go
    totalPrice := cart.CalculateTotalPrice()
    ```

7. Apply discounts and promotions:

    ```go
    cart.ApplyDiscount(discount)
    ```

8. Store cart data persistently:

    ```go
    cart.StoreData()
    ```
