# WORKER-REPORT-MATRIX

Cette lambda a pour but de générer le messsage contenant toutes les informations requises par la file SQS `{env}-pdf-matrix` qui invoque par la suite la lambda [worker-pdf-generator](https://github.com/ProovGroup/worker-pdf-generator) qui se charge de générer l'output matriciel du rapport.

# Fonctionnement

Cette lambda est invoqué par la file SQS `{env}-report-matrix`.

Le message reçu est sous le format suivant:
```json
{
    "proov_code": "XXXXXX",
}
```

La propriété `proov_code` permet de récupérer toutes les informations nécessaires afin de construire le message.

# Configuration

Le dossier `./assets` contient la structure de base du message qui est envoyé par cette lambda à la file `{env}-pdf-matrix`

Elle se construit de la manière suivante:

`./assets/message_structure.json`
```json
{
  "Source": {
    "Path": {
      "Region": "",
      "Bucket": "",
      "Key": ""
    },
    "Options": {
      "Type": "Archive"
    }
  },
  "Destination": {
    "Path": {
      "Region": "",
      "Bucket": "",
      "Key": ""
    },
    "Options": {
      "PDFOptions": {
        "Margin": {
          "Top": "20px",
          "Right": "50px",
          "Bottom": "100px",
          "Left": "50px"
        }
      }
    }
  },
  "DataSource": {
    "Content": null,
    "Options": {
      "Type": "Content"
    }
  }
}
```

Pour `Source` et `Destination` la propriété `Path` contient les chemins d'accès vers les dossier S3.

La propriété `Source` correspond au chemin vers le template utilisé afin de généré le matriciel.

La propriété `Destination` correspond au chemin de sortie du PDF généré.

La propriété `DataSource` contient les données du rapport dans `Content`, ce champs est remplie une fois les données récupérés et transformés.

> [!IMPORTANT]
> La configuration des chemins d'accès et d'autres variables de cette lambda se fait via `Configuration > Variables d'environnements` sur AWS.

