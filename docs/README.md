## Contents

We list a few examples of the `Image Generator`'s command here to illustrate its usefulness and ease of use.

- [Merge Command](#merge-command)
- [Result](#result)

## Merge Command

Use the `merge` command to merge several images from different directories into one.

```
Usage:
  generator merge [flags]

Flags:
  -d, --directory strings   specify a directory       
  -h, --help                help for merge
  -o, --out string          export the images
  -D, --root string         specify the root directory
```

let me try to explain how this tool work, please refer to the following diagram:

![Merge Command Diagram](https://www.plantuml.com/plantuml/png/bPDDJ-Cm48Rl_XMpFRLIcuY54oT4AG891ABT0n9IvSIJrE2Fo1yB4Th_lkC4-YKgn2dslA_7pqp6Ppvm5w9MS8kkHfXQhRaHS25vxRXclwGfLHG8gn0QVvhdZhzcjGohM4IFhhGce0bPQNNUA6vIfeaFGiaYdvvCxbFep-rDhhaFw2YqdD89BWswh24lOrVN5xFbV35xcDqj7kzdmn5ZvpOQmLqxB8_6C3ZdvKljKWWkhHBe1YDuPm3hHCEYNahDQ_JalkuM0rerfePAgdupRx2KLGjXXL1i4lR7eB8ED9KlJVradWpW6zUDmbCeLCGnE1HZQ54f-pe636Ks6B5_svY_4wOLxK50qdw6c99z1oTaLQ6ZxeD_aGPAgwpp-sZ7bFcIbtW8hIBjHsZfmBTeoRZ1e-4efE4m0MmHqdHf6yDg7_M4RhgBeOdX42mb-eJyJ1gaDNl0ezq2AKQoGnzLo9wzDScTlfqzRHP4s-k-Yq7Zx5yfg41fg4JTNGAy8EYkCc-ZeoROgk_33ih7RB1vImEl4t5wJdJ9plCY4tj0jJsv6oSNIU_axRdvoh3pzxuoh6Dzw_2jiIxiQnyjMejdQ0IzpVy0 "Merge Command Diagram")

You can download our [Example File](./assets/example.zip) for testing, there are five directories in the zip file, try to run the merge command below:

```
generator merge -D {YOUR_ROOT_DIRECTORY} -d background -d ears -d body -d eyes -d mouth -o {YOUR_OUTPUT_DIRECTORY}
```

As you can see, the merge command recognizes these options:

> Please note that if the field has been set, the field `-d` will adding the input as a prefix. for example: 
> ```
> generator -D `/root` -d `folder`
> ```
> It will look for the images starting from `/root/folder`

Name|Field|Type|Optional|Description|
----|-----|----|--------|-----------|
Root Direcoty|-D, --root|strings|true|specify the root directory, it is optional.
Direcoty|-d, --directory|string|false|specify the directory, you can reuse it.
Output|-o, --out|string|false|specify a directory for outputting `result.json` and images.
Help| -h, --help| |true|To see the usage document.

Then you will get the analysis of the combinations like this:

```
The analytical result:

background -> 3 
ears -> 3
body -> 3
eyes -> 3
mouth -> 3

Number of combinations: 243

Do you want to continue? (Y/N):
```

The merging program will run if you want to continue, and showing the errors:

```
Showing all errors below:

=== adding errors ===

=== merging errors ===

it's done, sucess count: 243
```

## Result

After the process of merging is done, you will get a json file for the result. As in the example below: 

```
[
  {
      "id": 241,
      "attributes": [
          {
              "trait_type": "background",
              "value": "bg-yellow"
          },
          {
              "trait_type": "ears",
              "value": "ear-tiger"
          },
          {
              "trait_type": "body",
              "value": "body-tiger"
          },
          {
              "trait_type": "eyes",
              "value": "eye-sleepy"
          },
          {
              "trait_type": "mouth",
              "value": "mouth-rat-19"
          }
      ],
      "path": "/usr/local/241.png"
  }
  ...
]
```