# Go Metadata Builder

A lightweight and efficient **Go-based metadata query builder** that enables dynamic schema definition, record updates, and metadata-driven operations. This tool simplifies metadata management by leveraging **JSON-based configurations** to define, validate, and execute database queries.

## 🚀 Key Features

- ✅ **JSON-Driven Schema & Queries** – Define table structures and operations dynamically.
- ✅ **Memory Optimization** – Designed with efficiency in mind for high-performance execution.
- ✅ **Custom Use Case Handling** – Supports various metadata-driven scenarios with extensibility.
- ✅ **Test Coverage & Robust Validation** – Ensures reliability with well-defined test cases.

## 📌 Use Cases

- Dynamic table schema management
- Metadata-driven query execution
- Scalable and flexible database operations

## 🔧 Installation

If using go modules.

```bash
go get -u github.com/nandrechetan/gomb
```

## 🛠 Features

1. Use the query builder to generate dynamic SQL queries.
2. Execute queries efficiently while ensuring data integrity.

## 🛠 Usage

```go
    table := gomb.NewTable("users")
    table.Comment = "Store user information"

    // Add columns
    table.AddColumn(gomb.NewColumn("id").SetDataType(gomb.SerialType).SetPrimaryKey().SetComment("User ID"))
    table.AddColumn(gomb.NewColumn("username").SetDataType(gomb.StringType).SetLength(50).SetNotNull().SetComment("Unique username"))
    table.AddColumn(gomb.NewColumn("email").SetDataType(gomb.StringType).SetLength(255).SetNotNull().SetComment("User email address"))
    table.AddColumn(gomb.NewColumn("password_hash").SetDataType(gomb.StringType).SetLength(100).SetNotNull())
    table.AddColumn(gomb.NewColumn("first_name").SetDataType(gomb.StringType).SetLength(50))
    table.AddColumn(gomb.NewColumn("last_name").SetDataType(gomb.StringType).SetLength(50))
    table.AddColumn(gomb.NewColumn("birth_date").SetDataType(gomb.DateType))
    table.AddColumn(gomb.NewColumn("is_active").SetDataType(gomb.BooleanType).SetDefault(gomb.DefaultTrue))
    table.AddColumn(gomb.NewColumn("login_count").SetDataType(gomb.IntegerType).SetDefault(0))
    table.AddColumn(gomb.NewColumn("last_login").SetDataType(gomb.DateTimeType))
    table.AddColumn(gomb.NewColumn("created_at").SetDataType(gomb.DateTimeType).SetDefault(gomb.DefaultCurrentTimestamp))
    table.AddColumn(gomb.NewColumn("updated_at").SetDataType(gomb.DateTimeType).SetDefault(gomb.DefaultCurrentTimestamp))

    sql, errors := table.ToSQL()

```

### Output:

"CREATE TABLE users (id SERIAL PRIMARY KEY COMMENT 'User ID', username VARCHAR(50) NOT NULL COMMENT 'Unique username', email VARCHAR(255) NOT NULL COMMENT 'User email address', password_hash VARCHAR(100) NOT NULL, first_name VARCHAR(50), last_name VARCHAR(50), birth_date DATE, is_active BOOLEAN DEFAULT TRUE, login_count INTEGER DEFAULT 0, last_login TIMESTAMP, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP) COMMENT ON TABLE users IS 'Store user information'"

## 🤝 Contributing

We welcome contributions! Please check out the issues tab for open tasks and improvements.

### How to contribute:

1. Fork the repository
2. Create a new branch (`git checkout -b feature-name`)
3. Commit your changes (`git commit -m "Added new feature"`)
4. Push to the branch (`git push origin feature-name`)
5. Create a pull request

## 📜 License

MIT License

Copyright (c) 2025 Nandre Chetan Kumar

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

🔗 **Stay Connected**\
📧 Contact: [nandrechetan@gmail.com](mailto:nandrechetan@gmail.com)\
🌐 GitHub: [nandrechetan](https://github.com/nandrechetan)
