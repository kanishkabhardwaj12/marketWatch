# 📊 MarketWatch

`MarketWatch` is a personal trade analysis and portfolio visualization tool that helps you understand your **Equity** and **Mutual Fund** positions through easy-to-navigate dashboards powered by **Grafana**.

It ingests tradebooks from Zerodha (Kite) and fetches real-time data from platforms like **MoneyControl** and **TickerTape**, then caches and aggregates your portfolio insights locally.

---

## ✨ Capabilities

- 📈 Visualize equity and mutual fund holdings using Grafana dashboards
- 📊 Track current market value of holdings
- 🔍 Visualise time of entry in comparison to the market trend
- 🔄 Caching to reduce API overhead

---

## 📸 Sample Dashboards

> _Screenshots go here_  
> *(Add images of Grafana panels showing top holdings, portfolio performance, and sector allocation)*

---

## 🚀 Getting Started

This guide walks you through the setup and usage of `marketWatch`.

---

### ✅ 1. Add a Configuration File

Create a file named `config.yaml` in the project root with the following format:

```yaml
mutual_funds:
  tradefiles_diretory: "./data/trade_books/MF"
equity:
  tradefiles_diretory: "./data/trade_books/EQ"
```  

---

### 🧾 2. Download Tradebook from Zerodha

To use `marketWatch`, you’ll need to download your historical tradebook files from [Zerodha Console](https://console.zerodha.com/). These are CSV files that contain your Equity and Mutual Fund transaction history.

#### 2.1 Steps to Download Tradebooks:

1. Log in to [https://console.zerodha.com/](https://console.zerodha.com/)
2. Navigate to **Reports** → **Downloads**
3. From the **Type** dropdown, select **Tradebook**
4. Choose the **Segment**:
    - _Equity_ for stocks (NSE/BSE)
    - _Mutual Funds_ for Coin transactions
5. Choose a suitable **Date Range** (preferably from your account inception)
6. Click **Download**
7. You should receive two CSV files:
    - `equity_tradebook.csv`
    - `mutual_fund_tradebook.csv`

Repeat the steps for both Equity and Mutual Fund segments if applicable.

---

### 📁 3. Place the Tradebook Files

Once you’ve downloaded the tradebook CSV files:

1. Create a folder (if not already) where tradebook files are expected to be stored.
2. Based on your `config.yaml`.

Note: Make sure they are all CSV Files

---

### ▶️ 4. Run the Tool

With Go installed on your system (version 1.18+ recommended):

You can run the service via:

```bash
make run
```

On startup, it will:

- Parse the tradebooks

- Cache relevant data from MoneyControl and TickerTape

- Start an HTTP server exposing the parsed and enriched data


---

### 📊 5. View the Dashboard

Once your server is up and running, you can explore your portfolio visually using the Grafana dashboard provided.

#### Steps:

1. Open your local Grafana instance in the browser (`http://localhost:3000` by default).
2. Log in (default username/password is usually `admin/admin` unless changed).
3. On the left sidebar, go to **Dashboards → Import**.
4. Click **Upload .json file** and select the file located at: `dashboard/grafana.json`
5. Assign a **data source** (this could be a JSON API, Prometheus, or any source your version of the dashboard is designed for).
6. Click **Import**.

Once imported, you’ll see panels showing:

- 🧾 Equity & MF holdings breakdown
- 📉 Current market value vs. invested amount
- 🧠 Sector & Fund House diversification
- ⏱️ Time-based entry visualization

you will require adding a datasource grafana-infinity-json
after that edit the pannels and choose the data source as the newly added one

---

## 🔐 Data Privacy
This is a fully offline tool:

- No trade or personal data is sent to any server.
- All processing is local.
- You remain in full control of your financial data.

---

## 📬 Contributions
Pull requests are welcome! If you have suggestions, open an issue or submit a PR. For major changes, please open a discussion first.

---


## 🙌 Acknowledgements
- Zerodha for providing exportable tradebooks
- TickerTape and MoneyControl for financial data
- Grafana for the awesome dashboards
