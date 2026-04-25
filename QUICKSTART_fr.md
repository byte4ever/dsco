# Guide de Démarrage Rapide dsco

Un tutoriel étape par étape pour maîtriser la gestion de configuration dsco.

## Table des Matières

1. [Installation](#1-installation)
2. [Concept Fondamental : Configuration par Pointeurs](#2-concept-fondamental--configuration-par-pointeurs)
3. [Votre Première Configuration](#3-votre-première-configuration)
4. [Comprendre les Couches](#4-comprendre-les-couches)
5. [Couches Struct : Valeurs par Défaut](#5-couches-struct--valeurs-par-défaut)
6. [Variables d'Environnement](#6-variables-denvironnement)
7. [Arguments en Ligne de Commande](#7-arguments-en-ligne-de-commande)
8. [Combiner Plusieurs Couches](#8-combiner-plusieurs-couches)
9. [Mode Strict : Sécurité de Configuration](#9-mode-strict--sécurité-de-configuration)
10. [Alias : Noms Raccourcis](#10-alias--noms-raccourcis)
11. [Fournisseurs Personnalisés](#11-fournisseurs-personnalisés)
12. [Gestion des Erreurs](#12-gestion-des-erreurs)
13. [Exemple Complet](#13-exemple-complet)

---

## 1. Installation

```bash
go get github.com/byte4ever/dsco
```

Nécessite Go 1.21 ou ultérieur.

---

## 2. Concept Fondamental : Configuration par Pointeurs

dsco utilise des **champs pointeurs** pour distinguer « non configuré » de
« configuré avec une valeur ». Cela évite les valeurs par défaut silencieuses
dangereuses.

### Le Problème avec la Configuration Traditionnelle

```go
// Approche traditionnelle - dangereuse !
type Config struct {
    Port    int    // 0 est-il intentionnel ou manquant ?
    Host    string // "" est-il intentionnel ou manquant ?
    Verbose bool   // false est-il intentionnel ou manquant ?
}
```

### La Solution dsco

```go
// Approche dsco - explicite et sûre
type Config struct {
    Port    *int    `yaml:"port"`    // nil = non configuré
    Host    *string `yaml:"host"`    // nil = non configuré
    Verbose *bool   `yaml:"verbose"` // nil = non configuré
}
```

**Point clé** : `nil` signifie « non configuré », toute valeur signifie
« explicitement défini ».

### L'Helper `R()`

dsco fournit `R[T](value T) *T` pour créer facilement des pointeurs :

```go
import "github.com/byte4ever/dsco"

// Au lieu de ceci :
port := 8080
config.Port = &port

// Utilisez ceci :
config.Port = dsco.R(8080)
```

---

## 3. Votre Première Configuration

Créons un exemple minimal fonctionnel :

```go
package main

import (
    "fmt"
    "log"

    "github.com/byte4ever/dsco"
)

type Config struct {
    Host *string `yaml:"host"`
    Port *int    `yaml:"port"`
}

func main() {
    var config *Config

    _, err := dsco.Fill(
        &config,
        dsco.WithStructLayer(&Config{
            Host: dsco.R("localhost"),
            Port: dsco.R(8080),
        }, "defaults"),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Serveur : %s:%d\n", *config.Host, *config.Port)
}
```

Sortie :
```
Serveur : localhost:8080
```

---

## 4. Comprendre les Couches

dsco utilise un **système de configuration en couches**. Chaque couche fournit
des valeurs, et **la première couche à fournir un champ l'emporte**. Les couches
suivantes ne remplissent que les champs laissés nil par toutes les couches
précédentes.

```
Couche 1 (première)  → priorité la plus haute (gagne)
Couche 2             → remplit les champs laissés nil par la Couche 1
Couche 3 (dernière)  → priorité la plus basse
```

Imaginez des transparents empilés — vous voyez le premier transparent
non-transparent en premier. Les couches inférieures n'apparaissent que là où
les couches supérieures sont transparentes (nil).

### Types de Couches

| Type de Couche | Source | Cas d'Usage |
|----------------|--------|-------------|
| `WithStructLayer` | Struct Go | Valeurs par défaut |
| `WithEnvLayer` | Variables d'environnement | Config conteneur/K8s |
| `WithCmdlineLayer` | Arguments ligne de commande | Surcharges à l'exécution |
| `WithStringValueProvider` | Fournisseur personnalisé | Secrets, fichiers, etc. |

---

## 5. Couches Struct : Valeurs par Défaut

Les couches struct fournissent des valeurs codées en dur, typiquement utilisées
pour les valeurs par défaut :

```go
type DatabaseConfig struct {
    Host    *string `yaml:"host"`
    Port    *int    `yaml:"port"`
    Timeout *int    `yaml:"timeout"`
}

defaults := &DatabaseConfig{
    Host:    dsco.R("localhost"),
    Port:    dsco.R(5432),
    Timeout: dsco.R(30),
}

_, err := dsco.Fill(
    &config,
    dsco.WithStructLayer(defaults, "defaults"),
)
```

Le second argument (`"defaults"`) est un identifiant pour les messages d'erreur.

### Valeurs par Défaut Partielles

Vous n'avez pas besoin de fournir tous les champs :

```go
// Ne fournir que certaines valeurs par défaut - les autres doivent venir d'autres couches
defaults := &DatabaseConfig{
    Port:    dsco.R(5432),  // Port par défaut
    Timeout: dsco.R(30),    // Timeout par défaut
    // Host intentionnellement nil - doit être fourni ailleurs
}
```

---

## 6. Variables d'Environnement

### Pourquoi les Préfixes sont Importants

Les préfixes de variables d'environnement sont essentiels en environnement de
production :

**1. Pods Multi-Conteneurs (Kubernetes)**

Quand plusieurs conteneurs s'exécutent dans un même pod, tous les conteneurs
partagent le même environnement. Les préfixes permettent de cibler la
configuration vers des conteneurs spécifiques :

```yaml
# Pod Kubernetes avec deux conteneurs
env:
  # Le conteneur frontend lit les variables FRONTEND-*
  - name: FRONTEND-PORT
    value: "8080"
  - name: FRONTEND-API-URL
    value: "http://localhost:3000"

  # Le conteneur backend lit les variables BACKEND-*
  - name: BACKEND-PORT
    value: "3000"
  - name: BACKEND-DATABASE-HOST
    value: "postgres.default.svc"
```

```go
// frontend/main.go
dsco.Fill(&config, dsco.WithEnvLayer("FRONTEND"))

// backend/main.go
dsco.Fill(&config, dsco.WithEnvLayer("BACKEND"))
```

**2. Éviter les Conflits**

Le système et les outils tiers définissent de nombreuses variables
d'environnement (`PATH`, `HOME`, `HTTP_PROXY`, `DATABASE_URL`, etc.). Les
préfixes empêchent votre application de lire accidentellement des variables
non liées ou d'entrer en conflit avec des variables existantes :

```bash
# Sans préfixe - risqué ! Pourrait entrer en conflit avec les vars système
HOST=localhost        # Pourrait entrer en conflit avec d'autres outils
PORT=8080             # Nom de variable courant

# Avec préfixe - sûr et explicite
MYAPP-HOST=localhost  # Appartient clairement à votre app
MYAPP-PORT=8080       # Pas d'ambiguïté
```

**3. Instances Multiples**

Exécutez plusieurs instances de la même application avec des configurations
différentes :

```bash
# Instance 1
WORKER1-QUEUE=high-priority WORKER1-CONCURRENCY=10 ./worker

# Instance 2
WORKER2-QUEUE=low-priority WORKER2-CONCURRENCY=5 ./worker
```

### Choisir de Bons Préfixes

**Évitez les préfixes génériques** - ils sont trop courants et créent de la
confusion :

```bash
# MAUVAIS : Préfixes génériques, ambigus
APP-HOST=...        # Quelle app ? Chaque service est une "app"
SERVER-PORT=...     # Quel serveur ? Trop vague
SERVICE-URL=...     # Sans signification dans un environnement microservices
CONFIG-TIMEOUT=...  # Tout a une config
```

**Utilisez des préfixes spécifiques basés sur le rôle** qui identifient le
composant :

```bash
# BON : Préfixes clairs et distinguables
API-HOST=...              # La passerelle API
WORKER-CONCURRENCY=...    # Worker de jobs en arrière-plan
CONSUMER-BATCH-SIZE=...   # Consommateur de file de messages
SCHEDULER-INTERVAL=...    # Service cron/planificateur
GATEWAY-RATE-LIMIT=...    # Passerelle API
INDEXER-CHUNK-SIZE=...    # Indexeur de recherche
NOTIFIER-SMTP-HOST=...    # Service de notification
```

**Pourquoi c'est important :**

1. **Débogage** : Quand vous voyez `WORKER-TIMEOUT=30` dans les logs, vous
   savez instantanément à quel service ça appartient

2. **Manifestes Kubernetes** : Des préfixes clairs rendent les fichiers YAML
   auto-documentés :
   ```yaml
   env:
     - name: ORDERAPI-DATABASE-HOST    # Évidemment pour l'API des commandes
     - name: PAYMENTWORKER-RETRY-MAX   # Évidemment pour le worker de paiement
   ```

3. **Environnements partagés** : En dev/staging où plusieurs services partagent
   l'infrastructure, des préfixes spécifiques évitent la contamination croisée
   accidentelle

4. **Communication d'équipe** : « Vérifie la config INDEXER » est plus clair
   que « vérifie la config APP pour le service indexeur »

**Conventions de nommage par type de service :**

| Type de Service | Exemples de Préfixes |
|-----------------|---------------------|
| APIs HTTP | `USERAPI`, `ORDERAPI`, `AUTHAPI` |
| Workers en arrière-plan | `EMAILWORKER`, `PAYMENTWORKER` |
| Consommateurs de messages | `ORDERCONSUMER`, `EVENTCONSUMER` |
| Jobs planifiés | `REPORTSCHEDULER`, `CLEANUPJOB` |
| Passerelles/proxies | `APIGATEWAY`, `AUTHPROXY` |

### Format Général

```
PREFIX-KEY=value
│      │
│      └─ Clé en MAJUSCULES (tirets et underscores autorisés)
└─ Préfixe (MAJUSCULES lettres et chiffres uniquement)
```

### Fonctionnement du Parsing

1. **Préfixe** : Doit correspondre à `^[A-Z][A-Z0-9]*$` (lettres/chiffres
   majuscules, commence par une lettre)
2. **Séparateur** : Un seul tiret (`-`) entre le préfixe et la clé
3. **Clé** : Doit correspondre à `^[A-Z][A-Z0-9]*(?:[-_][A-Z][A-Z0-9]*)*$`
   - Commence par une lettre majuscule
   - Peut contenir des tirets (`-`) ou underscores (`_`) comme séparateurs de mots
   - Chaque segment de mot commence par une lettre majuscule

### Correspondance Struct vers Variable d'Environnement

Étant donné cette struct :

```go
type Config struct {
    Host        *string         `yaml:"host"`
    Port        *int            `yaml:"port"`
    MaxRetry    *int            `yaml:"max_retry"`
    Database    *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host     *string `yaml:"host"`
    Port     *int    `yaml:"port"`
    PoolSize *int    `yaml:"pool_size"`
}
```

Avec `dsco.WithEnvLayer("MYAPP")`, la correspondance est :

| Chemin Struct | Clé YAML | Clé Interne | Variable d'Environnement |
|---------------|----------|-------------|-------------------------|
| `Config.Host` | `host` | `host` | `MYAPP-HOST` |
| `Config.Port` | `port` | `port` | `MYAPP-PORT` |
| `Config.MaxRetry` | `max_retry` | `max_retry` | `MYAPP-MAX_RETRY` |
| `Config.Database.Host` | `database.host` | `database-host` | `MYAPP-DATABASE-HOST` |
| `Config.Database.Port` | `database.port` | `database-port` | `MYAPP-DATABASE-PORT` |
| `Config.Database.PoolSize` | `database.pool_size` | `database-pool_size` | `MYAPP-DATABASE-POOL_SIZE` |

**Règles de transformation des clés :**
- Les chemins de structs imbriquées utilisent le **tiret** (`-`) comme séparateur de niveau
- Les noms de champs conservent leurs **underscores** (`_`) des tags yaml
- Tout après le préfixe est en **MAJUSCULES** dans la var env, **minuscules** en interne

### Exemples Valides vs Invalides

```bash
# Variables d'environnement valides
MYAPP-HOST=localhost           # Clé simple
MYAPP-MAX-RETRY=5              # Tiret dans la clé (imbriqué ou kebab-case)
MYAPP-DB_POOL_SIZE=10          # Underscore dans la clé (du tag yaml)
MYAPP-DATABASE-HOST=postgres   # Champ de struct imbriquée

# Variables d'environnement invalides
MYAPP_HOST=localhost           # Faux : underscore au lieu de tiret après le préfixe
myapp-HOST=localhost           # Faux : préfixe en minuscules
MYAPP-host=localhost           # Faux : clé en minuscules
MYAPP--HOST=localhost          # Faux : double tiret
MYAPP-123KEY=value             # Faux : clé commence par un chiffre
```

### Exemple d'Utilisation

```go
type Config struct {
    Host     *string         `yaml:"host"`
    Database *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host *string `yaml:"host"`
    Port *int    `yaml:"port"`
}

var config *Config

_, err := dsco.Fill(
    &config,
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(&Config{
        Host: dsco.R("localhost"),
        Database: &DatabaseConfig{
            Port: dsco.R(5432),
        },
    }, "defaults"),
)
```

Exécution :
```bash
MYAPP-HOST=api.example.com MYAPP-DATABASE-HOST=db.example.com MYAPP-DATABASE-PORT=5433 ./myapp
```

---

## 7. Arguments en Ligne de Commande

Les arguments en ligne de commande utilisent le format `--key=value` :

```go
_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)
```

Exécution :
```bash
./myapp --host=production.example.com --port=9000
```

### Règles de Format des Clés

Les clés doivent être en minuscules avec des tirets ou underscores :
- `--host=value` (clé simple)
- `--max-connections=100` (kebab-case)
- `--db_host=localhost` (snake_case)

**Invalide** : `--Host=value` (majuscules non autorisées)

### Champs Imbriqués

Pour les structs imbriquées, les clés sont jointes avec des **tirets** (pas
des points) :

```go
type Config struct {
    Database *DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host *string `yaml:"host"`
    Port *int    `yaml:"port"`
}
```

```bash
# Correct : séparé par des tirets
./myapp --database-host=db.example.com --database-port=5432

# Faux : les points ne sont PAS supportés
./myapp --database.host=db.example.com  # Invalide !
```

---

## 8. Combiner Plusieurs Couches

### Développement Local Rapide

Combiner des **couches struct** (valeurs par défaut) avec une **couche ligne
de commande** permet un développement local rapide sans aucune configuration
externe :

```go
var config *Config

_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),
    dsco.WithStructLayer(&Config{
        Host:     dsco.R("localhost"),
        Port:     dsco.R(8080),
        Database: &DatabaseConfig{
            Host: dsco.R("localhost"),
            Port: dsco.R(5432),
            Name: dsco.R("devdb"),
            User: dsco.R("devuser"),
        },
        LogLevel: dsco.R("debug"),
    }, "dev-defaults"),
)
```

**Avantages pour le développement local :**

1. **Zéro configuration pour démarrer** : Lancez simplement `./myapp` - toutes
   les valeurs par défaut sont intégrées au code, pas de fichiers de config ou
   de vars env nécessaires

2. **Surcharges rapides** : Testez différents scénarios sans éditer de fichiers :
   ```bash
   # Tester avec un port différent
   ./myapp --port=9000

   # Tester contre une base de données staging
   ./myapp --database-host=staging-db.example.com

   # Tester avec un logging de type production
   ./myapp --log-level=info
   ```

3. **Auto-documenté** : Les valeurs par défaut dans le code montrent quelles
   valeurs sont attendues et à quoi ressemble la configuration de développement

4. **Pas de pollution d'environnement** : Contrairement aux vars env, les
   arguments en ligne de commande ne persistent pas et n'affectent pas les
   autres processus

**Workflow typique :**
```bash
# Développement quotidien - lancez simplement
./myapp

# Test rapide avec un changement
./myapp --port=9000

# Tester un scénario spécifique
./myapp --database-host=testdb --log-level=trace

# Partager les étapes exactes de reproduction avec l'équipe
./myapp --feature-flag=true --timeout=5s
```

Ce pattern est particulièrement utile pour :
- Les microservices avec beaucoup d'options de configuration
- Les sessions de débogage rapides
- La reproduction de problèmes avec des paramètres spécifiques
- L'intégration de nouveaux développeurs (pas de setup requis)

### Pile de Couches Complète

La puissance de dsco vient de la combinaison de couches avec une priorité
claire :

```go
type Config struct {
    Host    *string `yaml:"host"`
    Port    *int    `yaml:"port"`
    Debug   *bool   `yaml:"debug"`
    Timeout *int    `yaml:"timeout"`
}

var config *Config

_, err := dsco.Fill(
    &config,
    // Couche 1 : Ligne de commande (priorité la plus haute)
    dsco.WithCmdlineLayer(),

    // Couche 2 : Variables d'environnement (priorité moyenne)
    dsco.WithEnvLayer("MYAPP"),

    // Couche 3 : Valeurs par défaut codées en dur (priorité la plus basse)
    dsco.WithStructLayer(&Config{
        Host:    dsco.R("localhost"),
        Port:    dsco.R(8080),
        Debug:   dsco.R(false),
        Timeout: dsco.R(30),
    }, "defaults"),
)
```

### Exemple de Priorité

Étant donné :
- Ligne de commande : `--host=production.example.com`
- Environnement : `MYAPP-HOST=staging.example.com`
- Couche struct : `Host="localhost"`, `Port=8080`

Résultat :
- `Host` = `"production.example.com"` (de cmdline — première couche à le fournir)
- `Port` = `8080` (de struct — cmdline et env l'ont laissé nil)

---

## 9. Mode Strict : Sécurité de Configuration

Le **mode strict** garantit que toutes les valeurs d'une couche stricte sont
consommées pendant le remplissage. Une valeur est « non consommée » si :

1. **Elle ne correspond à aucun champ de config** (détection de fautes de frappe)
2. **Elle a été remplacée par une couche antérieure** (détection de surcharge)

### Normal vs Strict

| Mode | Comportement |
|------|--------------|
| Normal | Valeurs non correspondantes et surcharges ignorées silencieusement |
| Strict | Valeurs non correspondantes et surcharges causent une erreur |

### Comprendre la Détection de Surcharge

Rappelez-vous : **la première couche à fournir un champ l'emporte**. Si la
valeur d'une couche stricte est remplacée par une couche antérieure dans la
liste, c'est une erreur.

```go
_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),            // Couche antérieure (position 0)
    dsco.WithStrictEnvLayer("MYAPP"),  // Couche stricte (position 1)
)
```

Si `--port` et `MYAPP-PORT` sont tous deux fournis, la valeur cmdline gagne
(couche antérieure). Mais comme la couche env est stricte, sa valeur remplacée
cause une `OverriddenKeyError`.

### Détection de Fautes de Frappe

Le mode strict détecte aussi les fautes de frappe - valeurs qui ne
correspondent à aucun champ de config :

```bash
# Faute de frappe : HOOST au lieu de HOST
MYAPP-HOOST=localhost ./myapp
# Erreur : "HOOST" ne correspond à aucun champ, reste inutilisé
```

### Le Positionnement Compte

Puisque la première couche à fournir un champ l'emporte :

- **Couche stricte en tête** → ses valeurs gagnent ; erreurs uniquement pour les champs non correspondants
- **Couche stricte en queue** → erreurs si des couches antérieures ont déjà fourni ses valeurs

```go
// Pattern 1 : Cmdline strict en tête (détection de fautes de frappe uniquement)
dsco.Fill(&config,
    dsco.WithStrictCmdlineLayer(),  // Erreurs uniquement pour les flags non correspondants
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)

// Pattern 2 : Env strict, s'assurer que les vars env ne sont pas remplacées par cmdline
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),
    dsco.WithStrictEnvLayer("MYAPP"),  // Erreurs si cmdline a déjà fourni le champ
)

// Pattern 3 : Valeurs immuables verrouillées en tête
dsco.Fill(&config,
    dsco.WithStrictStructLayer(&Config{
        APIEndpoint: dsco.R("https://api.production.com"),
    }, "immutable"),  // Priorité la plus haute, erreurs si non consommé
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithCmdlineLayer(),
)
```

### Quand Utiliser le Mode Strict

**Utilisez strict pour** :
- Garantir que certaines valeurs de couche ne peuvent pas être remplacées
- Détecter les fautes de frappe dans les noms de variables d'environnement
- Détecter les flags de ligne de commande inconnus
- S'assurer que les valeurs de configuration immuables sont utilisées

**Utilisez normal pour** :
- Les valeurs par défaut (peuvent être sautées si une couche antérieure fournit le champ)
- Les couches de priorité basse où être supplanté par des couches antérieures est attendu

---

## 10. Alias : Noms Raccourcis

Les alias fournissent des noms courts pour les clés de configuration :

```go
type Config struct {
    Database *DatabaseConfig `yaml:"database"`
    Server   *ServerConfig   `yaml:"server"`
    Logging  *LoggingConfig  `yaml:"logging"`
}

_, err := dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(
        dsco.WithAliases(map[string]string{
            // Format : "alias": "chemin.interne"
            "db-host": "database.host",   // --db-host → database-host
            "db-port": "database.port",   // --db-port → database-port
            "port":    "server.port",     // --port → server-port
            "v":       "logging.verbose", // --v → logging-verbose
        }),
    ),
    dsco.WithEnvLayer("MYAPP"),
    dsco.WithStructLayer(defaults, "defaults"),
)
```

Maintenant vous pouvez utiliser :
```bash
./myapp --db-host=localhost --port=9000 --v=true
```

Au lieu de :
```bash
./myapp --database-host=localhost --server-port=9000 --logging-verbose=true
```

**Note** : Le côté droit de la correspondance d'alias utilise des points
(format de chemin interne), mais les vraies clés de ligne de commande
utilisent des tirets.

---

## 11. Fournisseurs Personnalisés

Pour les sources de configuration au-delà de env/cmdline/structs, implémentez
un fournisseur personnalisé :

### L'Interface

```go
type NamedStringValuesProvider interface {
    GetName() string
    GetStringValues() svalue.Values
}
```

### Exemple : Fournisseur de Fichier

```go
import (
    "os"

    "github.com/byte4ever/dsco/svalue"
    "gopkg.in/yaml.v3"
)

type FileProvider struct {
    name   string
    values svalue.Values
}

func NewFileProvider(path string) (*FileProvider, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var raw map[string]string
    if err := yaml.Unmarshal(data, &raw); err != nil {
        return nil, err
    }

    values := make(svalue.Values)
    for k, v := range raw {
        values[k] = &svalue.Value{
            Value:    v,
            Location: svalue.NewLocation("file", path, k),
        }
    }

    return &FileProvider{name: path, values: values}, nil
}

func (f *FileProvider) GetName() string              { return f.name }
func (f *FileProvider) GetStringValues() svalue.Values { return f.values }
```

### Utilisation

```go
fileProvider, err := NewFileProvider("config.yaml")
if err != nil {
    log.Fatal(err)
}

_, err = dsco.Fill(
    &config,
    dsco.WithCmdlineLayer(),                     // CLI (priorité la plus haute)
    dsco.WithEnvLayer("MYAPP"),                  // Env
    dsco.WithStringValueProvider(fileProvider),  // Config fichier (priorité la plus basse)
)
```

### Exemple : Fournisseur de Secrets

```go
type VaultProvider struct {
    client *vault.Client
}

func (v *VaultProvider) GetName() string { return "vault" }

func (v *VaultProvider) GetStringValues() svalue.Values {
    values := make(svalue.Values)

    // Récupérer les secrets depuis Vault
    secret, _ := v.client.Read("secret/myapp")

    for k, val := range secret.Data {
        values[k] = &svalue.Value{
            Value:    val.(string),
            Location: svalue.NewLocation("vault", "secret/myapp", k),
        }
    }

    return values
}
```

---

## 12. Gestion des Erreurs

dsco fournit des erreurs détaillées avec suivi de localisation.

### Types d'Erreurs

| Type d'Erreur | Cause |
|---------------|-------|
| `LayerErrors` | Problèmes d'enregistrement de couche |
| `FillerErrors` | Problèmes de remplissage de struct |
| `InvalidInputError` | Type de cible invalide |
| `CmdlineAlreadyUsedError` | Plusieurs couches cmdline |
| `OverriddenKeyError` | Valeur stricte a été remplacée |

### Vérification des Erreurs

```go
_, err := dsco.Fill(&config, layers...)
if err != nil {
    var layerErr dsco.LayerErrors
    if errors.As(err, &layerErr) {
        for _, e := range layerErr.Errors() {
            log.Printf("Erreur de couche : %v", e)
        }
    }

    var fillerErr dsco.FillerErrors
    if errors.As(err, &fillerErr) {
        for _, e := range fillerErr.Errors() {
            log.Printf("Erreur de remplissage : %v", e)
        }
    }

    log.Fatal(err)
}
```

### Suivi de Localisation

dsco trace l'origine de chaque valeur :

```go
locations, err := dsco.Fill(&config, layers...)
if err != nil {
    log.Fatal(err)
}

// Afficher l'origine de chaque valeur
for path, loc := range locations {
    fmt.Printf("%s: %s\n", path, loc)
}
```

Sortie :
```
host: env[MYAPP-HOST]
port: cmdline[--port]
timeout: struct[defaults]
```

---

## 13. Exemple Complet

Voici un exemple prêt pour la production combinant tous les concepts :

### config.go

```go
package main

import (
    "errors"
    "time"

    "github.com/byte4ever/dsco"
)

// Config représente la configuration de l'application.
type Config struct {
    Server   *ServerConfig   `yaml:"server"`
    Database *DatabaseConfig `yaml:"database"`
    Logging  *LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
    Host         *string        `yaml:"host"`
    Port         *int           `yaml:"port"`
    ReadTimeout  *time.Duration `yaml:"read_timeout"`
    WriteTimeout *time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
    Host     *string `yaml:"host"`
    Port     *int    `yaml:"port"`
    Name     *string `yaml:"name"`
    User     *string `yaml:"user"`
    Password *string `yaml:"password"`
    SSLMode  *string `yaml:"ssl_mode"`
}

type LoggingConfig struct {
    Level   *string `yaml:"level"`
    Format  *string `yaml:"format"`
    Verbose *bool   `yaml:"verbose"`
}

// Validate vérifie les champs requis et les contraintes.
func (c *Config) Validate() error {
    if c.Server == nil || c.Server.Port == nil {
        return errors.New("server.port est requis")
    }
    if c.Database == nil || c.Database.Host == nil {
        return errors.New("database.host est requis")
    }
    if c.Database.Password == nil {
        return errors.New("database.password est requis")
    }
    return nil
}

// DefaultConfig retourne des valeurs par défaut sensées.
func DefaultConfig() *Config {
    return &Config{
        Server: &ServerConfig{
            Host:         dsco.R("0.0.0.0"),
            Port:         dsco.R(8080),
            ReadTimeout:  dsco.R(30 * time.Second),
            WriteTimeout: dsco.R(30 * time.Second),
        },
        Database: &DatabaseConfig{
            Port:    dsco.R(5432),
            SSLMode: dsco.R("require"),
        },
        Logging: &LoggingConfig{
            Level:   dsco.R("info"),
            Format:  dsco.R("json"),
            Verbose: dsco.R(false),
        },
    }
}
```

### main.go

```go
package main

import (
    "fmt"
    "log"

    "github.com/byte4ever/dsco"
)

func main() {
    var config *Config

    locations, err := dsco.Fill(
        &config,
        // 1. Ligne de commande (priorité la plus haute)
        dsco.WithCmdlineLayer(
            dsco.WithAliases(map[string]string{
                "db-host":     "database.host",
                "db-port":     "database.port",
                "db-name":     "database.name",
                "db-user":     "database.user",
                "db-password": "database.password",
                "port":        "server.port",
                "v":           "logging.verbose",
            }),
        ),

        // 2. Variables d'environnement
        dsco.WithStrictEnvLayer("APP"),

        // 3. Valeurs par défaut (priorité la plus basse)
        dsco.WithStructLayer(DefaultConfig(), "defaults"),
    )
    if err != nil {
        log.Fatalf("Erreur de configuration : %v", err)
    }

    // Valider les champs requis
    if err := config.Validate(); err != nil {
        log.Fatalf("Erreur de validation : %v", err)
    }

    // Afficher les sources de configuration
    fmt.Println("Configuration chargée depuis :")
    for path, loc := range locations {
        fmt.Printf("  %s: %s\n", path, loc)
    }

    // Démarrer l'application
    fmt.Printf("\nDémarrage du serveur sur %s:%d\n",
        *config.Server.Host,
        *config.Server.Port,
    )
}
```

### Exécution de l'Exemple

```bash
# Avec les valeurs par défaut uniquement (échouera à la validation - pas de mot de passe db)
./myapp

# Avec les valeurs requises
APP-DATABASE-HOST=db.example.com \
APP-DATABASE-USER=appuser \
APP-DATABASE-PASSWORD=secret123 \
APP-DATABASE-NAME=mydb \
./myapp

# Surcharger le port via la ligne de commande
APP-DATABASE-HOST=db.example.com \
APP-DATABASE-USER=appuser \
APP-DATABASE-PASSWORD=secret123 \
./myapp --port=9000 --db-name=production -v=true
```

---

## Résumé

| Concept | Point Clé |
|---------|-----------|
| Champs pointeurs | `nil` = non configuré, valeur = explicitement défini |
| Helper `R()` | Crée des pointeurs facilement : `dsco.R(8080)` |
| Priorité des couches | La première couche à fournir un champ l'emporte |
| Couches struct | Valeurs par défaut codées en dur |
| Variables env | Format : `PREFIX-KEY=value` |
| Ligne de commande | Format : `--key=value` |
| Mode strict | Erreurs sur les valeurs non utilisées |
| Alias | Noms courts pour les chemins imbriqués |
| Fournisseurs personnalisés | Implémenter `NamedStringValuesProvider` |

**Bonne Pratique** : Toujours ordonner les couches de la priorité la plus haute
à la plus basse ; la première couche à fournir un champ l'emporte :
```go
dsco.Fill(&config,
    dsco.WithCmdlineLayer(),        // 1. Ligne de commande (priorité la plus haute)
    dsco.WithEnvLayer(...),         // 2. Environnement
    dsco.WithStringValueProvider(), // 3. Fichiers/Secrets
    dsco.WithStructLayer(...),      // 4. Valeurs par défaut (priorité la plus basse)
)
```

---

## Quelles clés dois-je fournir ?

`inventory.Compute` liste toutes les clés que votre configuration en couches
attend, sans rien lire dans l'environnement, les arguments ou les fichiers :

```go
report, _ := inventory.Compute(&config, layers...)
report.WriteText(os.Stdout)
```

Voir la [section Inventaire du README](README_fr.md#inventaire) pour la
description complète et les trois exemples exécutables (sortie texte, JSON
pour l'outillage, et une vérification préalable qui fait échouer la CI
lorsque des clés requises manquent).
