#!/bin/bash

echo "🔧 Test des touches de navigation - watch-fs"
echo "=========================================="
echo ""

echo "📋 Instructions de test :"
echo "1. L'application va se lancer dans 3 secondes"
echo "2. Vous devriez voir les événements des fichiers créés"
echo "3. Testez les touches suivantes :"
echo ""
echo "   Navigation :"
echo "   - ↑/↓/←/→ (flèches)"
echo "   - h/j/k/l (vim-style)"
echo "   - Page Up/Page Down"
echo "   - u/d (page alternative)"
echo "   - Home/End"
echo "   - g/G (début/fin)"
echo ""
echo "   Autres :"
echo "   - f : Basculer fichiers"
echo "   - d : Basculer répertoires"
echo "   - a : Basculer agrégation"
echo "   - s : Changer tri"
echo "   - q : Quitter"
echo ""

echo "⏳ Lancement dans 3 secondes..."
sleep 3

echo "🚀 Lancement de watch-fs..."
./bin/watch-fs -path . 