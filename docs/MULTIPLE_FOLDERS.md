# Surveillance de plusieurs dossiers

Cette fonctionnalité permet de surveiller plusieurs dossiers simultanément avec `watch-fs`.

## Utilisation

### Nouveau flag `-paths`

Utilisez le flag `-paths` pour spécifier plusieurs dossiers à surveiller, séparés par des virgules :

```bash
# Surveiller 3 dossiers
./watch-fs -paths "/path/to/dir1,/path/to/dir2,/path/to/dir3"

# Surveiller 2 dossiers avec espaces
./watch-fs -paths "/path/to/dir1, /path/to/dir2"

# Surveiller plusieurs dossiers avec des chemins relatifs
./watch-fs -paths "./src,./tests,./docs"
```

### Compatibilité avec l'ancien flag `-path`

Le flag `-path` original reste fonctionnel pour la compatibilité avec les versions antérieures :

```bash
# Ancienne syntaxe (toujours supportée)
./watch-fs -path "/single/directory"
```

## Fonctionnalités

### Affichage dans l'interface

L'interface TUI affiche maintenant tous les dossiers surveillés :

- **Un seul dossier** : Affiche le chemin complet
- **2-3 dossiers** : Affiche tous les chemins séparés par des virgules
- **Plus de 3 dossiers** : Affiche les 2 premiers + le nombre de dossiers restants

Exemples d'affichage :

```
Watching: /path/to/dir1 | Events: 5 | Sort: Time
Watching: /path/to/dir1, /path/to/dir2, /path/to/dir3 | Events: 12 | Sort: Time
Watching: /path/to/dir1, /path/to/dir2, +3 more (5 dirs) | Events: 8 | Sort: Time
```

### Gestion des erreurs

- Si un dossier n'existe pas, une erreur est affichée et le programme s'arrête
- Si aucun flag de chemin n'est fourni, un message d'aide s'affiche
- Les espaces autour des virgules sont automatiquement supprimés

### Validation

Tous les dossiers spécifiés sont validés avant de démarrer la surveillance :

- Vérification de l'existence des dossiers
- Vérification que ce sont bien des dossiers (pas des fichiers)

## Exemples d'utilisation

### Développement web

```bash
# Surveiller les dossiers de développement
./watch-fs -paths "./src,./public,./assets"
```

### Surveillance système

```bash
# Surveiller plusieurs dossiers système
./watch-fs -paths "/var/log,/tmp,/home/user/documents"
```

### Projet multi-modules

```bash
# Surveiller différents modules d'un projet
./watch-fs -paths "./frontend/src,./backend/src,./shared"
```

## Tests

La fonctionnalité est entièrement testée avec :

- Tests unitaires pour la logique de surveillance multiple
- Tests d'intégration pour les flags de ligne de commande
- Tests de compatibilité avec l'ancienne syntaxe
- Tests de gestion d'erreurs

Exécutez les tests avec :

```bash
./test/test_multiple_folders.sh
```

## Migration depuis l'ancienne version

Si vous utilisez actuellement `-path`, vous pouvez continuer à l'utiliser. Pour migrer vers la nouvelle syntaxe :

**Avant :**

```bash
./watch-fs -path "/my/directory"
```

**Après (optionnel) :**

```bash
./watch-fs -paths "/my/directory"
```

Pour ajouter des dossiers supplémentaires, utilisez simplement `-paths` avec plusieurs chemins.
