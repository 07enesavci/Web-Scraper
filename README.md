# Go Web Scraper & Screenshot Tool

Bu proje, Go (Golang) kullanılarak geliştirilmiş, modern ve performanslı bir web kazıma (scraping) aracıdır. `Chromedp` kütüphanesi kullanılarak gerçek bir tarayıcı (headless chrome) üzerinde çalışır.

## 🚀 Özellikler

Bu araç, belirtilen bir web sitesini ziyaret eder ve aşağıdaki işlemleri otomatik olarak gerçekleştirir:

-   **Tam Sayfa Ekran Görüntüsü :** Sitenin en üstünden en altına kadar yüksek kaliteli (Full HD) ekran görüntüsü alır.
-   **HTML Kaynak Kodu İndirme:** Sayfanın işlenmiş (render edilmiş) son HTML halini kaydeder.
-   **Link Çıkarma:** Sayfadaki tüm bağlantıları (`<a>` etiketleri) toplar ve bir listeye yazar.
-   **Hata Yönetimi:**
    -   HTTP Status kodlarını (200, 404, 500) kontrol eder.
    -   404 (Sayfa Bulunamadı) veya 500 (Sunucu Hatası) durumlarında gereksiz dosya oluşturmaz, kullanıcıyı uyarır.
-   **Zaman Aşımı (Timeout) Koruması:** 30 saniye içinde yanıt vermeyen sitelerde işlem otomatik olarak iptal edilir.
-   **Gürültü Önleme:** Gereksiz tarayıcı loglarını gizler, temiz bir çıktı sunar.

## 🛠️ Kurulum

Bu projeyi çalıştırmak için bilgisayarınızda [Go](https://go.dev/dl/) yüklü olmalıdır.

1.  Projeyi klonlayın:
    ```bash
    git clone https://github.com/07enesavci/Web-Scraper.git
    cd Web-Scraper
    ```

2.  Gerekli kütüphaneleri indirin:
    ```bash
    go mod download
    ```

## 💻 Kullanım

Programı çalıştırmak için terminalden `go run` komutunu ve hedef web sitesini kullanmanız yeterlidir.

**Örnek 1 (https ile):**
```bash
go run main.go https://yildizcti.com/
