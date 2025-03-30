describe('Create a new animal', () => {
  beforeEach(() => {
    cy.intercept('GET', 'http://localhost:3000/animals', { fixture: 'animals.json' }).as('getAnimals')
    cy.visit('/')
  })

  it('open modal, fill form and add the new animal to the table', () => {
    cy.get('button').contains('New').click()

    cy.get('.modal').should('be.visible')

    cy.get('input[id="name"]').type('Simba')
    cy.get('input[id="species"]').type('Lion')
    cy.get('input[id="birthday"]').type('2020-01-01')
    cy.get('input[id="genre"]').type('male')
    cy.get('input[id="diet"]').type('Carnivore')
    cy.get('input[id="condition"]').type('Healthy')
    cy.get('textarea[id="notes"]').type('Young lion')

    cy.get('button').contains('Add').click()

    cy.get('.modal').should('not.exist')

    // cy.get('table').contains('Simba').should('exist')
  })
})

