#!/bin/bash

# Test Chat History - Simulasi percakapan dengan konteks

echo "🧪 Test 1: Pesan pertama - perkenalan"
echo "=================================="
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Halo, nama saya Taufan",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7819,
      "speed": 0,
      "traffic": "smooth"
    },
    "history": []
  }'

echo -e "\n\n"
read -p "Tekan Enter untuk lanjut ke test 2..."

echo "🧪 Test 2: Pesan kedua - tanya nama (dengan history)"
echo "=================================="
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Siapa nama saya?",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7819,
      "speed": 0,
      "traffic": "smooth"
    },
    "history": [
      {
        "role": "user",
        "content": "Halo, nama saya Taufan"
      },
      {
        "role": "assistant",
        "content": "Halo Taufan! Senang berkenalan dengan Anda. Saya adalah asisten polisi lalu lintas AI yang siap membantu Anda dengan informasi lalu lintas, rute, dan keselamatan berkendara. Ada yang bisa saya bantu hari ini?"
      }
    ]
  }'

echo -e "\n\n"
read -p "Tekan Enter untuk lanjut ke test 3..."

echo "🧪 Test 3: Skenario Samsat - pesan pertama"
echo "=================================="
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Saya mau ke Kantor Samsat Tangsel",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7819,
      "speed": 0,
      "traffic": "smooth"
    },
    "history": []
  }'

echo -e "\n\n"
read -p "Tekan Enter untuk lanjut ke test 4..."

echo "🧪 Test 4: Skenario Samsat - follow-up (dengan history)"
echo "=================================="
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Berapa jaraknya dari lokasi saya?",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7819,
      "speed": 0,
      "traffic": "smooth"
    },
    "history": [
      {
        "role": "user",
        "content": "Saya mau ke Kantor Samsat Tangsel"
      },
      {
        "role": "assistant",
        "content": "Baik! Kantor Samsat Tangsel (Tangerang Selatan) berlokasi di Jl. Pajajaran No.100, Pamulang, Kota Tangerang Selatan. Dari lokasi Anda saat ini di Jakarta Selatan, saya akan bantu cari rute terbaik untuk kesana."
      }
    ]
  }'

echo -e "\n\n✅ Test selesai!"
