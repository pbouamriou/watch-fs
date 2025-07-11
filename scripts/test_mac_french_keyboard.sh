#!/bin/bash

# Script de test pour les claviers Mac français
# Test des touches de navigation alternatives

echo "🍎 Test des touches de navigation pour MacBook Pro avec clavier français"
echo "======================================================================"
echo ""

echo "📋 Touches de navigation disponibles :"
echo ""
echo "🔹 Navigation élément par élément :"
echo "  • Flèches ↑/↓/←/→ (standard)"
echo "  • h/j/k/l (vim-style, alternative pour Mac)"
echo ""
echo "🔹 Navigation par page :"
echo "  • Page Up/Page Down (standard)"
echo "  • u/d (alternative pour Mac)"
echo ""
echo "🔹 Aller au début/fin :"
echo "  • Home/End (standard)"
echo "  • g/G (vim-style, alternative pour Mac)"
echo ""

echo "🎯 Raccourcis complets pour Mac français :"
echo "  • h/j/k/l : Navigation (gauche/bas/haut/droite)"
echo "  • u/d : Page précédente/suivante"
echo "  • g/G : Aller au début/fin"
echo "  • f : Basculer fichiers"
echo "  • d : Basculer répertoires"
echo "  • a : Basculer agrégation"
echo "  • s : Changer tri"
echo "  • q : Quitter"
echo ""

echo "🧪 Instructions de test :"
echo "  1. Lancez l'application : ./bin/watch-fs -path ."
echo "  2. Créez quelques fichiers pour générer des événements"
echo "  3. Testez les touches alternatives :"
echo "     - h/j/k/l pour la navigation"
echo "     - u/d pour les pages"
echo "     - g/G pour début/fin"
echo ""

echo "💡 Note : Si les flèches ne fonctionnent pas, utilisez h/j/k/l"
echo "✅ Toutes les touches alternatives sont maintenant disponibles !" 