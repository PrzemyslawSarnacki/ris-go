![CI Status](https://github.com/finalrep/ris-go/actions/workflows/ci.yml/badge.svg)
[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)


# ris-go – Relative Index for Streetlifting Classic

This is modification of the ris-go library to use it with SL Classic competitions. Main modification is calculation of Classic total instead of All4 and adding constraints to the existing model so it matches Classic needs. Heavier athletes lift more absolute weight, but strength doesn't scale linearly with bodyweight. With SL Classic athletes lift less overall weight without Squats and Muscle-ups, so we recalibrated the target benchmark down from ~600 kg totals to ~240 kg totals but The curve remains intact. Below are calculated parameters as of date 2nd October 2026: 
```
=== 2-LIFT STREETLIFTING PARAMETERS ===

MALE PARAMETERS:
A: 111.00000
K: 237.05671
B: 0.12507
v: 70.24735
Q: 0.36688
RMSE: 23.14 kg

FEMALE PARAMETERS:
A: 45.00000
K: 96.00000
B: 0.11089
v: 81.38343
Q: 0.01361
RMSE: 12.45 kg
```

Male Classic Excel Formula
Put A2 = Bodyweight (kg) and B2 = Pull-up + Dip Total (kg):

```
=B2*100/(111 + (237.05671 - 111)/(1 + 0.36688*EXP(-0.12507*(A2 - 70.24735))))
```

Female Classic Excel Formula
Put A2 = Bodyweight (kg) and B2 = Pull-up + Dip Total (kg):

```
=B2*100/(45 + (96 - 45)/(1 + 0.01361*EXP(-0.11089*(A2 - 81.38343))))
```

**ris-go** is a Go library for calculating the *Relative Index for Streetlifting (RIS)*. The library allows the calculation of the RIS value based on individual strength levels and body weight, and it provides a method for determining the parameters of the underlying mathematical model based on real sports data.

This project was [initially created](https://warisradji.com/ris/) by [Waris Radji](https://www.instagram.com/riiswa) and [Mathieu Ardoin](https://www.instagram.com/mat95.sw) as part of a bachelor's thesis with the goal of making performance evaluation methods in strength sports comparable and adaptable.

---

## Background

The RIS is a relative scoring system specifically for the sport of **Streetlifting**. It is designed to make athletes of different weight classes comparable by normalizing total performance.

### Formula

$$
\text{RIS} = \frac{\text{Total} \times 100}{A + \frac{K - A}{1 + Q \cdot e^{-B \cdot (\text{BW} - v)}}}
$$

**Parameters:**

- `Total`: Total performance (e.g., sum of Weighted Pull-Up and Weighted Dip)
- `BW`: Body weight
- `A, K, Q, B, v`: Parameters that are optimized through fitting to real data

---

## Features

- Calculation of the RIS value with given parameters
- Fitting of RIS parameters to performance data via nonlinear optimization
- Importing and processing CSV data
- Modular architecture in Go

---

## Quick Start

### Prerequisites

- Go ≥ 1.24
- [`gonum`](https://github.com/gonum/gonum) for mathematical optimization

### Installation

```bash
go get github.com/finalrep/ris-go/lib
```

#### (Optional) Python setup

only needed if you want to use `FitRISParamsScipy`

```python
python3 -m venv .venv  # Creation of a virtualenv is recommanded
source .venv/bin/activate
pip install -r requirements.txt
```

### Usage

We recommend to use `FitRISParamsNelder` which uses NelderMead for fitting optimization. The original implementation can be used by using `FitRISParamsScipy` but needs extra setup for python.

### Example

- see our [example](https://github.com/FinalRep/ris-go/tree/main/examples) implementation
- The following plots use the 2023 data and our example to create a fitting for both male and female athletes

![data plot male](https://github.com/finalrep/ris-go/blob/main/examples/male.png?raw=true)
![data plot female](https://github.com/finalrep/ris-go/blob/main/examples/female.png?raw=true)

## Development and Contribution

- _Test your code properly!_
- use [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/)
- use [semantic versioning](https://semver.org/lang/de/) to create tags and versions
- use [pre-commit](https://pre-commit.com/)

### pre-commit hook

This project uses [pre-commit](https://pre-commit.com/).
Install pre-commit and run•
```
pre-commit install
```

### Using mage

Install [Mage](https://github.com/magefile/mage/) by

```
go install github.com/magefile/mage@latest
mage -init

# outputs all possible usage options
mage
```

### Testing

Run `mage test:run` to test the code and `mage test:cover` to see test coverage
