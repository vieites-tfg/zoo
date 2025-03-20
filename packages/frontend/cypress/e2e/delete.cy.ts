describe('Delete an animal', () => {
  beforeEach(() => {
    cy.intercept('GET', 'http://localhost:3000/animals', { fixture: 'animals.json' }).as('getAnimals')
    cy.visit('/')
    cy.wait('@getAnimals')
  })

  it('selects an animal and deletes it', () => {
    cy.get('table tbody tr').contains('Leo').parents('tr').within(() => {
      cy.get('input[type="checkbox"]').check()
    })

    cy.get('button').contains('Delete')
      .should('not.have.class', 'cursor-not-allowed')
      .click()

    // cy.get('table').contains('Leo').should('not.exist')
  })
})

