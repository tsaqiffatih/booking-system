Untuk **batas waktu pembayaran** (payment expiry time) dalam transaksi Midtrans, ya, Anda bisa mengatur batas waktu pembayaran agar transaksi hanya berlaku dalam jangka waktu tertentu, seperti **1 jam** setelah pesanan dibuat. Namun, fitur ini tidak langsung diatur dalam API **CoreAPI** ataupun **Snap**. Sebagai gantinya, Anda dapat mengaturnya dengan cara berikut:

### 1. **Mengatur Expiry Time pada Transaksi**

Midtrans memungkinkan Anda untuk menetapkan batas waktu untuk setiap transaksi dengan menggunakan field **expiry** pada objek `TransactionDetails` saat membuat transaksi menggunakan **CoreAPI** atau **Snap**. Expiry ini menentukan kapan transaksi akan kadaluarsa.

#### **Contoh Pengaturan Expiry Time menggunakan CoreAPI**

Saat Anda membuat transaksi, Anda bisa menentukan waktu **expiry** dalam format waktu UNIX atau dengan menggunakan **`expiry`** sebagai bagian dari transaksi:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veritrans/go-midtrans"
)

func main() {
	// Inisialisasi dengan ServerKey dan ClientKey Anda
	midtrans.ServerKey = "your-server-key"
	midtrans.ClientKey = "your-client-key"

	// Menentukan waktu kadaluarsa (misalnya 1 jam dari sekarang)
	expiryTime := time.Now().Add(1 * time.Hour) // 1 jam dari waktu sekarang
	expiryUnix := expiryTime.Unix()              // Mengonversi waktu ke format UNIX timestamp

	// Membuat transaksi
	transaction := &midtrans.TransactionDetails{
		OrderID:  "ORDER123456",
		GrossAmt: 100000, // Total pembayaran (dalam sen)
	}

	item := midtrans.ItemDetails{
		ID:    "ITEM123",
		Name:  "Product Example",
		Price: 100000,
		Qty:   1,
	}

	// Membuat charge request dengan expiry
	chargeRequest := &midtrans.ChargeReq{
		TransactionDetails: *transaction,
		ItemDetails:        []midtrans.ItemDetails{item},
		PaymentType:        midtrans.PaymentTypeCreditCard, // Metode pembayaran (bisa disesuaikan)
		Expiry:             &midtrans.Expiry{ // Menetapkan waktu kadaluarsa
			Duration: 60,  // Durasi dalam menit (60 menit = 1 jam)
			Unit:     midtrans.ExpiryUnitMinute,
		},
	}

	// Melakukan charge transaksi
	chargeResponse, err := midtrans.ChargeTransaction(chargeRequest)
	if err != nil {
		log.Fatalf("Error in charging transaction: %v", err)
	}

	// Menampilkan hasil transaksi
	fmt.Printf("Transaction Status: %s\n", chargeResponse.Status)
	fmt.Printf("Transaction ID: %s\n", chargeResponse.TransactionID)
	fmt.Printf("Redirect URL: %s\n", chargeResponse.RedirectURL) // URL untuk menyelesaikan pembayaran

	// 2. Mengecek Status Transaksi
	statusResponse, err := midtrans.CheckTransaction("ORDER123456") // Ganti dengan OrderID yang relevan
	if err != nil {
		log.Fatalf("Error in checking transaction status: %v", err)
	}

	// Menampilkan status transaksi
	fmt.Println("Transaction Status:", statusResponse.Status)
}
```

### Penjelasan:
- **`Expiry`**: Anda mengatur durasi kedaluwarsa transaksi dengan menambahkan waktu tertentu (`Duration` dan `Unit`). Dalam contoh ini, kita menggunakan durasi **60 menit** (1 jam). Midtrans akan otomatis menandai transaksi sebagai kedaluwarsa setelah waktu tersebut.
- **`ExpiryUnitMinute`**: Anda bisa memilih unit waktu seperti menit, jam, atau hari sesuai kebutuhan. Pada contoh di atas, kita menggunakan **menit**.

### 2. **Mengatur Pembayaran untuk Transaksi yang Kedaluwarsa**

Jika batas waktu untuk pembayaran sudah terlewat, transaksi akan **berstatus expired**. Anda dapat memeriksa status transaksi menggunakan fungsi `CheckTransaction`. Jika statusnya **expired**, maka Anda bisa meminta pengguna untuk membuat pesanan baru.

#### **Contoh Memeriksa Status Transaksi**

```go
// Mengecek status transaksi
statusResponse, err := midtrans.CheckTransaction("ORDER123456") // Ganti dengan OrderID yang relevan
if err != nil {
	log.Fatalf("Error in checking transaction status: %v", err)
}

// Memeriksa apakah transaksi sudah expired
if statusResponse.Status == "expired" {
	fmt.Println("Transaction expired. Please make a new order.")
} else {
	fmt.Println("Transaction Status:", statusResponse.Status)
}
```

Jika transaksi sudah kedaluwarsa (misalnya lebih dari 1 jam), status transaksi akan menjadi **expired**. Pada titik ini, Anda dapat menginstruksikan pengguna untuk membuat pesanan baru.

### 3. **Mengatur Pembayaran Menggunakan Snap**

Jika Anda menggunakan **Snap**, Anda juga bisa mengatur **expiry time** pada transaksi saat membuat Snap Token. Ini memungkinkan Anda untuk menetapkan batas waktu pembayaran bagi pengguna.

```go
snapRequest := &midtrans.SnapReq{
	TransactionDetails: midtrans.TransactionDetails{
		OrderID:  "ORDER123456",
		GrossAmt: 100000,
	},
	ItemDetails: []midtrans.ItemDetails{
		{
			ID:    "ITEM123",
			Name:  "Product Example",
			Price: 100000,
			Qty:   1,
		},
	},
	Expiry: &midtrans.Expiry{
		Duration: 60, // Durasi dalam menit (60 menit = 1 jam)
		Unit:     midtrans.ExpiryUnitMinute,
	},
}

// Membuat Snap Token untuk transaksi ini
snapTokenResponse, err := midtrans.Snap.CreateTransaction(snapRequest)
if err != nil {
	log.Fatalf("Error creating Snap transaction: %v", err)
}

// Menampilkan token dan URL pembayaran
fmt.Println("Snap Token:", snapTokenResponse.Token)
fmt.Println("Payment URL:", snapTokenResponse.RedirectURL)
```

### 4. **Mengelola Transaksi Kedaluwarsa**

Jika transaksi kedaluwarsa, Anda bisa menangani logika di aplikasi Anda:
- **Jika transaksi kedaluwarsa**: Tampilkan pesan kepada pengguna untuk melakukan pemesanan baru dan ulangi proses pembayaran.
- **Jika transaksi berhasil**: Lanjutkan dengan proses pengiriman produk atau layanan sesuai dengan bisnis Anda.

### 5. **Kesimpulan**

- **Batas waktu pembayaran** bisa diatur dengan menentukan **`expiry`** pada objek transaksi. Anda dapat menetapkan waktu kadaluarsa dalam bentuk **menit**, **jam**, atau **hari**.
- **Status transaksi** akan menjadi **expired** setelah batas waktu tersebut berlalu, dan Anda bisa menangani status ini dengan memeriksa **`CheckTransaction`** atau menggunakan logika aplikasi untuk meminta pengguna membuat pesanan baru.
- Penggunaan **Snap** dan **CoreAPI** mendukung fitur ini, jadi Anda bisa menggunakan cara yang sesuai dengan integrasi Anda.

Dengan begitu, Anda bisa membatasi waktu pembayaran sesuai dengan kebutuhan bisnis, seperti misalnya transaksi hanya berlaku selama 1 jam setelah pemesanan dibuat.