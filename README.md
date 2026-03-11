# gophish_fileAttachmentModding

A hardened GoPhish fork for authorized phishing simulations and security awareness campaigns.
Un fork reforzado de GoPhish para simulaciones de phishing autorizadas y campañas de concienciación.

---

## Modifications / Modificaciones

| Change | Detail |
|---|---|
| `X-Gophish-Contact` removed | Eliminado de cabeceras de correo |
| `X-Gophish-Signature` renamed | → `X-Hub-Signature-256` |
| `X-Mailer` spoofed | Simula Microsoft Outlook 16.0 |
| `X-Server: gophish` removed | Eliminado del servidor HTTP |
| `Server` header | Responde como `Apache/2.4.54 (Ubuntu)` |
| 404 page | Página de error estilo Apache |
| Attachment tracking | Nuevas rutas `/download` y `/track-open` |
| `gendoc` CLI | Genera DOCX/XLSX con pixel de tracking embebido |

---

## Requirements / Requisitos

- Go 1.18+
- GCC (for sqlite3)

```bash
sudo apt-get install -y golang-go gcc
```

---

## Build / Compilar

```bash
git clone https://github.com/Juguitos/gophish_fileAttachmentModding.git
cd gophish_fileAttachmentModding

# Main server / Servidor principal
go build -ldflags "-X github.com/gophish/gophish/config.Version=0.12.1" -o cyberphish .

# Document generator / Generador de documentos
go build -o gendoc ./cmd/gendoc/
```

---

## Run / Ejecutar

```bash
./cyberphish
```

- Admin panel: `https://<your-domain>:3333`
- Phishing server: `http://<your-domain>:80`

---

## Attachment Campaigns / Campañas con Adjunto

### 1. Generate tracking document / Generar documento de tracking

```bash
# DOCX
./gendoc -url "http://your-domain.com/track-open?rid=PLACEHOLDER" \
         -o static/attachments/Invoice_March.docx

# XLSX
./gendoc -url "http://your-domain.com/track-open?rid=PLACEHOLDER" \
         -o static/attachments/Report_March.xlsx
```

### 2. Email template / Plantilla de correo

```html
<a href="http://your-domain.com/download?rid={{.RId}}&f=Invoice_March.docx">
  Download Invoice
</a>
```

> **Important:** Use your domain directly, **not** `{{.URL}}`.
> `{{.URL}}` already appends `?rid=...` and will break the link.

### 3. Tracked events / Eventos registrados

| Event | Trigger |
|---|---|
| `Downloaded Attachment` | User clicks the download link |
| `Opened Attachment` | User opens the file in Word/Excel (embedded pixel) |

---

## config.json example

```json
{
  "admin_server": {
    "listen_url": "0.0.0.0:3333",
    "use_tls": true,
    "cert_path": "gophish_admin.crt",
    "key_path": "gophish_admin.key",
    "trusted_origins": []
  },
  "phish_server": {
    "listen_url": "0.0.0.0:80",
    "use_tls": false,
    "cert_path": "example.crt",
    "key_path": "example.key"
  },
  "db_name": "sqlite3",
  "db_path": "gophish.db",
  "migrations_prefix": "db/db_",
  "contact_address": "",
  "logging": {
    "filename": "",
    "level": ""
  }
}
```

---

## Notes / Notas

- The embedded tracking pixel in DOCX/XLSX requires Word/Excel to load external content (may be blocked by corporate policies).
- El pixel en DOCX/XLSX requiere que Word/Excel cargue contenido externo (puede estar bloqueado en entornos corporativos).
- Download tracking (`/download`) always works regardless of Office settings.
- El tracking de descarga (`/download`) funciona siempre independientemente de la configuración de Office.

---

**For authorized security awareness testing only. / Solo para pruebas de concienciación de seguridad autorizadas.**
