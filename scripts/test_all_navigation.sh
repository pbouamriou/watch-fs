#!/bin/bash

echo "🧪 Test complet de la navigation - watch-fs"
echo "=========================================="
echo ""

echo "📋 Création d'événements de test..."
for i in {1..15}; do
    echo "test event $i" > "nav_test_$i.txt"
done

echo "✅ 15 fichiers créés pour tester la navigation"
echo ""

echo "🎯 Test à effectuer :"
echo "1. Lancez l'application : ./bin/watch-fs -path ."
echo "2. Vous devriez voir 15 événements CREATE"
echo "3. Testez toutes les touches de navigation :"
echo ""
echo "   Navigation élément par élément :"
echo "   - ↑/↓ : Déplacement ligne par ligne"
echo "   - h/j/k/l : Navigation vim-style"
echo ""
echo "   Navigation par page :"
echo "   - Page Up/Page Down : Déplacement par page (10 lignes)"
echo "   - u/d : Navigation par page alternative"
echo ""
echo "   Navigation début/fin :"
echo "   - Home/End : Aller au début/fin"
echo "   - g/G : Navigation début/fin alternative"
echo ""
echo "   Autres touches :"
echo "   - f : Basculer fichiers"
echo "   - d : Basculer répertoires"
echo "   - a : Basculer agrégation"
echo "   - s : Changer tri"
echo "   - q : Quitter"
echo ""

echo "🔍 Points à vérifier :"
echo "  ✅ Le curseur se déplace d'une seule ligne à chaque pression"
echo "  ✅ Le curseur ne sort jamais de la liste"
echo "  ✅ Home/End et g/G fonctionnent correctement"
echo "  ✅ Page Up/Page Down fonctionnent correctement"
echo ""

echo "�� Prêt à tester !" 