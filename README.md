# Devoir Maison : Audit de Sécurité de l'Application Flop-Security

Les scripts développés pour l'exploitation des failles sont présents dans ce dépôt :
-   **`BruteForce/bruteforce.go`** : Un outil en ligne de commande développé en Go pour faire une attaque par force brute sur la page de connexion.
-   **`Selenium/attack-automatiser.py`** : Un script Python qui utilise Selenium pour automatiser une attaque en chaîne : connexion par injection SQL, puis injection d'un code XSS dans la section des commentaires.

---

## Injection SQL

### Description détaillée
L'injection SQL est une attaque qui consiste à insérer des commandes SQL malveillantes dans les requêtes envoyées par une application à sa base de données. L'application `flopsecurity` s'est révélée vulnérable via son formulaire de connexion.

En entarnt le champ `email` par ceci, il a été possible de contourner l'authentification :
```sql
" OR "1"="1" #
````

Cette entrée modifie la requête SQL exécutée par le serveur pour qu'elle soit toujours vraie, accordant ainsi l'accès sans mot de passe valide.

### Risques sur le système

  - **Confidentialité :** Vol de toutes les données de la base (informations utilisateurs, mots de passe, etc.).
  - **Intégrité :** Modification ou suppression de données.
  - **Disponibilité :** Un attaquant pourrait rendre le site inutilisable en supprimant des tables.
  - **Compromission du serveur :** Dans les cas les plus graves, peut mener à une exécution de commandes sur le serveur lui-même.

### Sécurisation

1.  **Requêtes Préparées :** C'est la défense la plus efficace. Il faut séparer les commandes SQL des données fournies par l'utilisateur pour que ces dernières ne soient jamais interprétées comme du code. En PHP, cela se fait via PDO ou MySQLi.
2.  **Principe du Moindre Privilège :** La stratégie de la base de données du TP est un bon exemple. L'utilisateur applicatif `flopsecurity` n'a de droits que sur sa propre base, limitant l'impact en cas de compromission.
3.  **Validation des Données :** Toujours vérifier et nettoyer les entrées utilisateur pour s'assurer qu'elles correspondent au format attendu.

-----

## XSS (Cross-Site Scripting)

### Description détaillée

L'attaque XSS consiste à injecter un script côté client dans une page web, qui sera ensuite exécuté par les navigateurs des autres utilisateurs. Dans ce TP, une faille **XSS stockée** a été exploitée.

Le script `Selenium/attack-automatiser.py` injecte le payload suivant dans la section des commentaires :

```html
<h2>Cette page est vulnérable !</h2>
```

Ce code, une fois enregistré et affiché, modifie le contenu de la page pour tous les visiteurs.

### Risques sur le système

Les risques ciblent principalement les utilisateurs de l'application :

  - **Vol de session :** Le risque majeur, permettant à un attaquant de voler le cookie de session d'un utilisateur et d'usurper son identité.
  - **Défiguration du site (Defacement) :** Modification du contenu visuel du site.
  - **Phishing et redirection :** Rediriger les utilisateurs vers des sites malveillants pour voler leurs identifiants.

### Sécurisation

1.  **Échappement en Sortie (Output Encoding) :** C'est la règle d'or. Avant d'afficher une donnée provenant d'un utilisateur, il faut systématiquement "neutraliser" les caractères HTML spéciaux. En PHP, la fonction `htmlspecialchars()` est conçue pour cela.
2.  **Content Security Policy (CSP) :** Un en-tête HTTP qui définit les sources de contenu (scripts, styles) approuvées, empêchant le navigateur de charger des scripts depuis des sources inconnues.
3.  **Validation en Entrée :** Refuser les entrées contenant du code potentiellement dangereux si le champ n'est pas censé en avoir.

-----

## Brute Force

### Description détaillée

Une attaque par force brute consiste à essayer systématiquement toutes les combinaisons possibles de mots de passe pour un utilisateur donné jusqu'à trouver la bonne.

Une démonstration de cette attaque a été réalisée avec le script `BruteForce/bruteforce.go`, qui lit une liste de mots de passe depuis `passwords.txt` et les teste via des requêtes HTTP POST.

### Risques sur le système

  - **Prise de contrôle de comptes :** Accès non autorisé aux comptes utilisateurs et à leurs données.
  - **Déni de Service (DoS) :** Le grand nombre de tentatives de connexion peut surcharger le serveur et le rendre inaccessible.

### Sécurisation

  - **Limitation des tentatives :** Bloquer un compte ou une adresse IP après N échecs.
  - **Délais progressifs (Throttling) :** Augmenter le temps d'attente entre chaque tentative échouée.
  - **CAPTCHA :** Demander une vérification humaine après plusieurs échecs.
  - **Politique de mots de passe robustes :** Imposer des mots de passe longs et complexes.
  - **Authentification Multi-Facteurs (MFA) :** La défense la plus efficace, requérant une seconde preuve d'identité.

-----

## Sécurisation d'une base MySQL

La stratégie de gestion des utilisateurs de la base de données présentée dans le TP `flopsecurity` améliore la sécurité en appliquant le **principe du moindre privilège**.

Elle segmente les accès en trois rôles distincts :

1.  **`root@localhost`** : Le super-utilisateur, dont l'accès est restreint à la machine locale uniquement. Cela empêche les tentatives de connexion `root` à distance, qui sont une cible majeure.
2.  **`dba@'%'`** : Un compte administrateur avec tous les privilèges, destiné aux tâches de maintenance (sauvegardes, modifications de structure). Il est distinct de `root` et son mot de passe doit être très robuste.
3.  **`flopsecurity@localhost`** : L'utilisateur applicatif. Son rôle est le plus restreint : il ne peut accéder **qu'à la base de données `flopsecurity`** et uniquement depuis la machine locale. En cas de compromission de l'application via une injection SQL, les dégâts sont confinés à cette seule base de données. L'attaquant ne peut ni voir, ni affecter les autres bases du serveur.

Cette séparation est fondamentale pour contenir l'impact d'une faille de sécurité.

-----

## Connexion SSH

Pour sécuriser `openssh-server`, il faut modifier le fichier `/etc/ssh/sshd_config` avec les paramètres suivants :

```sshd_config
# Interdire la connexion directe de l'utilisateur root.
PermitRootLogin no

# Désactiver l'authentification par mot de passe pour forcer l'utilisation de clés SSH, bien plus sécurisées.
PasswordAuthentication no
PubkeyAuthentication yes

# N'autoriser que les utilisateurs spécifiés (ou les membres d'un groupe spécifique).
AllowUsers uruz paul

# Utiliser uniquement le protocole SSHv2.
Protocol 2

# Changer le port par défaut pour éviter les scans automatisés.
# Port 2222
```

Après modification, le service doit être rechargé avec `sudo systemctl reload sshd`.

-----

## Firewall

### Rôle d'un firewall

Un pare-feu (firewall) agit comme une barrière filtrante entre un réseau (ou une machine) et le monde extérieur (Internet). Son rôle est de contrôler le trafic réseau entrant et sortant en se basant sur un ensemble de règles de sécurité. Il n'autorise que les communications qui ont été explicitement permises.

### Configuration

Dans le cadre du TP, le pare-feu `ufw` a été utilisé. Sa configuration est simple :

1.  `sudo ufw enable` : Pour l'activer.
2.  `sudo ufw default deny incoming` : Pour bloquer toutes les connexions entrantes par défaut (principe de sécurité).
3.  `sudo ufw allow ssh` ou `sudo ufw allow 22/tcp` : Pour autoriser des services ou des ports spécifiques.

On peut configurer des règles basées sur le port, le protocole (TCP/UDP), l'adresse IP source ou destination, et l'action (ALLOW/DENY).

-----

## Rôle des autres solutions de sécurisation

  - **Mises à jour système :** C'est une mesure de sécurité **critique**. Les mises à jour (via `apt update && apt upgrade`) corrigent les failles de sécurité découvertes dans le système d'exploitation et les logiciels installés (Apache, MySQL, etc.). Un système non mis à jour est une porte d'entrée facile pour un attaquant.

  - **Configuration sécurisée d'apache2 :** Cela inclut la désactivation de modules non nécessaires, la restriction des informations affichées sur les pages d'erreur, et le contrôle des permissions sur les dossiers. La **ré-écriture d'URL (`mod_rewrite`)** n'est pas une fonctionnalité de sécurité en soi, mais elle y contribue en masquant la structure réelle des fichiers (`/user/123` au lieu de `/get_user.php?id=123`), ce qui rend la reconnaissance plus difficile pour un attaquant.

-----

## Veille technologique

### Publications de l'ANSSI

L'ANSSI (Agence Nationale de la Sécurité des Systèmes d'Information) fournit des guides de référence pour la sécurisation :

  - **Le Guide d'hygiène informatique :** Contient les 42 règles essentielles pour sécuriser un système d'information.
  - **Recommandations pour le développement d'applications web sécurisées :** Un guide détaillé pour les développeurs sur la manière d'éviter les failles communes.
  - **Guides de configuration de systèmes :** Des recommandations spécifiques pour durcir la configuration de systèmes GNU/Linux, serveurs web, etc.

### Risques identifiés par OWASP

L'OWASP Top 10 liste les risques les plus critiques pour la sécurité des applications web. Voici un résumé des principaux risques de 2021 :

1.  **Broken Access Control :** Des utilisateurs peuvent accéder à des fonctions ou des données sans en avoir les droits.
2.  **Cryptographic Failures :** Des données sensibles sont mal protégées, souvent stockées en clair ou avec des algorithmes de chiffrement faibles.
3.  **Injection :** Les failles d'injection (SQL, NoSQL, etc.) restent une menace majeure.
4.  **Insecure Design :** La sécurité n'est pas intégrée dès la phase de conception de l'application.
5.  **Security Misconfiguration :** Des configurations par défaut non sécurisées, des messages d'erreur trop détaillés, etc.
6.  **Vulnerable and Outdated Components :** Utilisation de librairies, frameworks ou dépendances avec des failles de sécurité connues.
7.  **Identification and Authentication Failures :** Mauvaise gestion des sessions, mots de passe faibles, absence de protection contre le brute force.
8.  **Software and Data Integrity Failures :** Manque de vérification de l'intégrité des mises à jour ou des données, menant à l'installation de code malveillant.
