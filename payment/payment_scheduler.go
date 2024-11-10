package payment_service

import (
	"context"
	"log"
	"time"
)

// untuk merubah status dari payment saat user belum melakukan pembayaran.
func StartPaymentExpiryScheduler(repo PaymentRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-ticker.C
			if err := repo.CancelExpiredPayments(context.Background()); err != nil {
				log.Println("Error cancelling expired payments:", err)
			} else {
				log.Println("Successfully cancelled expired payments.")
			}
		}
	}()
}
