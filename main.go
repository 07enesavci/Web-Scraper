package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network" // ağ trafiği için gerekli kütüphane
	"github.com/chromedp/chromedp"        // tarayıcıdan işlem yapmak için gerekli kütüphane
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Hata: Bir web sitesi adresi girmeyi unuttun.\nKullanım: go run main.go https://google.com")
	}

	hedefAdres := os.Args[1]
	fmt.Println("🌐 İşlem başlıyor, gidilecek adres:", hedefAdres)

	tarayiciAyarlari := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1920, 1080),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	yoneticiContext, iptalYonetici := chromedp.NewExecAllocator(context.Background(), tarayiciAyarlari...)
	defer iptalYonetici()

	tarayiciContext, iptalTarayici := chromedp.NewContext(yoneticiContext,
		chromedp.WithLogf(func(string, ...interface{}) {}),
	)
	defer iptalTarayici()

	// 30 saniye zaman aşımı
	zamanIsleyici, iptalZaman := context.WithTimeout(tarayiciContext, 30*time.Second)
	defer iptalZaman()

	var htmlKaynakKod string
	var ekranGoruntusu []byte
	var linkListesi []string
	var httpDurumKodu int64

	// ağı dinleme
	chromedp.ListenTarget(zamanIsleyici, func(olay interface{}) {
		if gelenVeri, tamam := olay.(*network.EventResponseReceived); tamam {
			if gelenVeri.Type == network.ResourceTypeDocument {
				httpDurumKodu = gelenVeri.Response.Status
			}
		}
	})

	// yapılacak işlemler
	hata := chromedp.Run(zamanIsleyici,
		network.Enable(),
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(hedefAdres),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML(`html`, &htmlKaynakKod, chromedp.ByQuery),
		chromedp.FullScreenshot(&ekranGoruntusu, 90),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &linkListesi),
	)

	// hata analizi
	if hata != nil {
		hataMesaji := hata.Error()

		if strings.Contains(hataMesaji, "ERR_NAME_NOT_RESOLVED") {
			log.Fatalf("❌ HATA: Böyle bir site bulunamadı. Adresi yanlış yazmış olabilirsiniz.")
		} else if strings.Contains(hataMesaji, "ERR_CONNECTION_REFUSED") {
			log.Fatalf("❌ HATA: Site bağlantıyı reddetti. Sunucu kapalı olabilir.")
		} else if strings.Contains(hataMesaji, "context deadline exceeded") {
			log.Fatalf("❌ HATA: Site 30 saniye içinde açılmadı (Zaman Aşımı).")
		} else {
			// Bilinmeyen başka bir hata varsa teknik detayını göster
			log.Fatalf("❌ HATA: Siteye bağlanılamadı. Teknik Detay: %v", hata)
		}
	}

	// http durum kontrolu
	if httpDurumKodu != 0 && httpDurumKodu >= 400 {
		switch httpDurumKodu {
		case 404:
			log.Fatalf("❌ HATA: Sayfa Bulunamadı (404). Site var ama girdiğiniz sayfa yok.")
		case 500:
			log.Fatalf("❌ HATA: Sunucu Hatası (500). Karşı tarafın sistemi bozuk.")
		case 403:
			log.Fatalf("❌ HATA: Erişim Reddedildi (403). Bu siteye girmeniz yasak.")
		default:
			log.Fatalf("❌ HATA: Site hata kodu döndürdü: %d", httpDurumKodu)
		}
	}

	fmt.Printf("✅ Bağlantı başarılı! (HTTP Durumu: %d)\n", httpDurumKodu)

	// dosyaya kaydetme
	if hata := os.WriteFile("site_data.html", []byte(htmlKaynakKod), 0644); hata != nil {
		log.Fatal("HTML kaydedilemedi:", hata)
	}

	if hata := os.WriteFile("screenshot.png", ekranGoruntusu, 0644); hata != nil {
		log.Fatal("Resim kaydedilemedi:", hata)
	}

	linkMetni := strings.Join(linkListesi, "\n")
	if hata := os.WriteFile("links.txt", []byte(linkMetni), 0644); hata != nil {
		log.Fatal("Linkler kaydedilemedi:", hata)
	}

	fmt.Println("📸 Ekran görüntüsü alındı: screenshot.png")
	fmt.Println("📄 Kaynak kodlar kaydedildi: site_data.html")
	fmt.Printf("🔗 Toplam %d link bulundu: links.txt\n", len(linkListesi))
}
